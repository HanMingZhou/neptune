import fs from 'node:fs'
import path from 'node:path'

const rootDir = path.resolve(process.cwd(), 'src')
const allowedExtensions = new Set(['.ts', '.tsx', '.vue'])
const violations = []

const visit = (currentPath) => {
  const stats = fs.statSync(currentPath)
  if (stats.isDirectory()) {
    for (const entry of fs.readdirSync(currentPath)) {
      visit(path.join(currentPath, entry))
    }
    return
  }

  if (!allowedExtensions.has(path.extname(currentPath))) {
    return
  }

  const content = fs.readFileSync(currentPath, 'utf8')
  const lines = content.split(/\r?\n/)

  lines.forEach((line, index) => {
    if (/\bconsole\.log\s*\(/.test(line)) {
      violations.push(
        `${path.relative(process.cwd(), currentPath)}:${index + 1}:${line.trim()}`
      )
    }
  })
}

visit(rootDir)

if (violations.length > 0) {
  console.error('Unexpected console.log statements found:')
  violations.forEach((violation) => console.error(`- ${violation}`))
  process.exit(1)
}

console.log('No console.log statements found in src/.')
