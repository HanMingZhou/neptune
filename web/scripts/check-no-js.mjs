import fs from 'node:fs'
import path from 'node:path'

const projectRoot = process.cwd()
const sourceRoot = path.join(projectRoot, 'src')

const args = process.argv.slice(2)
const strictMode = args.includes('--strict')
const allowFileIndex = args.indexOf('--allow-file')
const allowFilePath =
  allowFileIndex >= 0 && args[allowFileIndex + 1]
    ? path.resolve(projectRoot, args[allowFileIndex + 1])
    : null

const normalizePath = (filePath) => filePath.split(path.sep).join('/')

const walk = (directory) => {
  const entries = fs.readdirSync(directory, { withFileTypes: true })
  const files = []

  for (const entry of entries) {
    const fullPath = path.join(directory, entry.name)
    if (entry.isDirectory()) {
      files.push(...walk(fullPath))
      continue
    }

    if (entry.isFile() && fullPath.endsWith('.js')) {
      files.push(normalizePath(path.relative(projectRoot, fullPath)))
    }
  }

  return files
}

const jsFiles = walk(sourceRoot).sort()

if (strictMode) {
  if (jsFiles.length > 0) {
    console.error(
      'Found JavaScript files in src/. Remove them before enabling strict no-js mode:'
    )
    jsFiles.forEach((filePath) => console.error(`- ${filePath}`))
    process.exit(1)
  }

  console.log('No JavaScript files found in src/.')
  process.exit(0)
}

if (!allowFilePath) {
  console.error('Missing --allow-file <path> argument.')
  process.exit(1)
}

if (!fs.existsSync(allowFilePath)) {
  console.error(`Allowlist file not found: ${allowFilePath}`)
  process.exit(1)
}

const allowList = JSON.parse(fs.readFileSync(allowFilePath, 'utf8'))
const allowSet = new Set(allowList.map(normalizePath))

const unexpectedFiles = jsFiles.filter((filePath) => !allowSet.has(filePath))
const missingFromSource = allowList
  .map(normalizePath)
  .filter((filePath) => !jsFiles.includes(filePath))

if (unexpectedFiles.length > 0) {
  console.error(
    'Detected new JavaScript files outside the approved migration allowlist:'
  )
  unexpectedFiles.forEach((filePath) => console.error(`- ${filePath}`))
  process.exit(1)
}

if (missingFromSource.length > 0) {
  console.log('Allowlist entries already migrated and ready to be removed:')
  missingFromSource.forEach((filePath) => console.log(`- ${filePath}`))
}

console.log(
  `JavaScript migration guard passed. Remaining allowlisted JS files: ${jsFiles.length}.`
)
