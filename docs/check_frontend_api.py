 import os
import re

def get_backend_apis(api_file="server/source/system/api.go"):
    apis = []
    with open(api_file, "r", encoding="utf-8") as f:
        content = f.read()
        # {ApiGroup: "...", Method: "...", Path: "...", Description: "..."}
        matches = re.finditer(r'\{ApiGroup:\s*"([^"]+)",\s*Method:\s*"([^"]+)",\s*Path:\s*"([^"]+)"', content)
        for m in matches:
            apis.append({
                "group": m.group(1),
                "method": m.group(2).lower(),
                "path": m.group(3)
            })
    return apis

def get_frontend_paths(api_dir="web/src/api"):
    paths = set()
    for root, dirs, files in os.walk(api_dir):
        for f in files:
            if f.endswith(".ts"):
                with open(os.path.join(root, f), "r", encoding="utf-8") as file:
                    content = file.read()
                    # url: '...' or url: "..."
                    matches = re.findall(r'url:\s*[\'"]([^\'"]+)[\'"]', content)
                    for path in matches:
                        paths.add(path)
    return paths

if __name__ == "__main__":
    backend_apis = get_backend_apis()
    frontend_paths = get_frontend_paths()
    
    missing = []
    for api in backend_apis:
        if api["path"] not in frontend_paths:
            missing.append(api)
            
    print(f"Total backend APIs: {len(backend_apis)}")
    print(f"Total frontend paths: {len(frontend_paths)}")
    print("\nMissing in Frontend:")
    for m in missing:
        # Exclude those mentioned by user
        if any(x in m["path"] for x in ["/init/", "/autocode/"]):
            continue
        print(f'{m["method"].upper()} {m["path"]} (Group: {m["group"]})')
