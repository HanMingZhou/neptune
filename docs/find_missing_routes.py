import os
import re

def extract_routes(folder="../server/router"):
    routes = []
    for root, dirs, files in os.walk(folder):
        for f in files:
            if f.endswith(".go"):
                path = os.path.join(root, f)
                with open(path, "r", encoding="utf-8") as file:
                    content = file.read()
                    
                    groups = {}
                    for match in re.finditer(r'([a-zA-Z0-9_]+)\s*:=\s*[A-Za-z0-9_]+\.Group\("([^"]+)"\)', content):
                        groups[match.group(1)] = match.group(2)
                    
                    for match in re.finditer(r'([a-zA-Z0-9_]+)\s*\.\s*Group\("([^"]+)"\)', content):
                        groups[match.group(1)] = match.group(2)
                        
                    for match in re.finditer(r'([a-zA-Z0-9_]+)\.(GET|POST|PUT|DELETE|Any)\("([^"]+)"[^\n]*?//\s*(.*)', content):
                        var_name = match.group(1)
                        method = match.group(2)
                        route_path = match.group(3)
                        desc = match.group(4).strip()
                        
                        group_path = ""
                        for g in groups:
                            if g == var_name:
                                group_path = groups[g]
                        
                        routes.append({
                            "file": path,
                            "var": var_name,
                            "method": method,
                            "path": route_path,
                            "group": group_path,
                            "desc": desc
                        })

                    for match in re.finditer(r'([a-zA-Z0-9_]+)\.(GET|POST|PUT|DELETE|Any)\("([^"]+)"', content):
                        var_name = match.group(1)
                        method = match.group(2)
                        route_path = match.group(3)
                        
                        existing = [r for r in routes if r['file'] == path and r['method'] == method and r['path'] == route_path]
                        if not existing:
                            group_path = ""
                            for g in groups:
                                if g == var_name:
                                    group_path = groups[g]
                                    
                            routes.append({
                                "file": path,
                                "var": var_name,
                                "method": method,
                                "path": route_path,
                                "group": group_path,
                                "desc": ""
                            })

    return routes

if __name__ == "__main__":
    with open("../server/source/system/api.go", "r", encoding="utf-8") as f:
        api_content = f.read()
    
    existing_apis = []
    # match existing recorded api paths
    for match in re.finditer(r'Path:\s*"([^"]+)"', api_content):
        existing_apis.append(match.group(1))
        
    for r in extract_routes():
        url_path = f'/api/v1/{r["group"] + "/" if r["group"] else ""}{r["path"]}'.replace('//', '/')
        if url_path.endswith('/'):
            url_path = url_path[:-1]
        combo = url_path
        if combo not in existing_apis:
            print(f'Missing: {r["method"]} {url_path} ({r["file"]}) -> var {r["var"]}, desc: {r["desc"]}')
