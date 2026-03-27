package mcpTool

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/model/system"
	"gin-vue-admin/service"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

// 注册工具
func init() {
	RegisterTool(&MenuCreator{})
}

// MenuCreateRequest 菜单创建请求结构
type MenuCreateRequest struct {
	ParentId    uint                   `json:"parentId"`    // 父菜单ID，0表示根菜单
	Path        string                 `json:"path"`        // 路由path
	Name        string                 `json:"name"`        // 路由name
	Hidden      bool                   `json:"hidden"`      // 是否在列表隐藏
	Component   string                 `json:"component"`   // 对应前端文件路径
	Sort        int                    `json:"sort"`        // 排序标记
	Title       string                 `json:"title"`       // 菜单名
	Icon        string                 `json:"icon"`        // 菜单图标
	KeepAlive   bool                   `json:"keepAlive"`   // 是否缓存
	DefaultMenu bool                   `json:"defaultMenu"` // 是否是基础路由
	CloseTab    bool                   `json:"closeTab"`    // 自动关闭tab
	ActiveName  string                 `json:"activeName"`  // 高亮菜单
	Parameters  []MenuParameterRequest `json:"parameters"`  // 路由参数
	MenuBtn     []MenuButtonRequest    `json:"menuBtn"`     // 菜单按钮
}

// MenuParameterRequest 菜单参数请求结构
type MenuParameterRequest struct {
	Type  string `json:"type"`  // 参数类型：params或query
	Key   string `json:"key"`   // 参数key
	Value string `json:"value"` // 参数值
}

// MenuButtonRequest 菜单按钮请求结构
type MenuButtonRequest struct {
	Name string `json:"name"` // 按钮名称
	Desc string `json:"desc"` // 按钮描述
}

// MenuCreateResponse 菜单创建响应结构
type MenuCreateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	MenuID  uint   `json:"menuId"`
	Name    string `json:"name"`
	Path    string `json:"path"`
}

// MenuCreator 菜单创建工具
type MenuCreator struct{}

// New 创建菜单创建工具
func (m *MenuCreator) New() mcp.Tool {
	return mcp.NewTool("create_menu",
		mcp.WithDescription(`创建前端菜单记录，用于AI编辑器自动添加前端页面时自动创建对应的菜单项。

**重要限制：**
- 当使用gva_auto_generate工具且needCreatedModules=true时，模块创建会自动生成菜单项，不应调用此工具
- 仅在以下情况使用：1) 单独创建菜单（不涉及模块创建）；2) AI编辑器自动添加前端页面时`),
		mcp.WithNumber("parentId",
			mcp.Description("父菜单ID，0表示根菜单"),
			mcp.DefaultNumber(0),
		),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("路由path，如：userList"),
		),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("路由name，用于Vue Router，如：userList"),
		),
		mcp.WithBoolean("hidden",
			mcp.Description("是否在菜单列表中隐藏"),
		),
		mcp.WithString("component",
			mcp.Required(),
			mcp.Description("对应的前端Vue组件路径，如：view/user/list.vue"),
		),
		mcp.WithNumber("sort",
			mcp.Description("菜单排序号，数字越小越靠前"),
			mcp.DefaultNumber(1),
		),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("菜单显示标题"),
		),
		mcp.WithString("icon",
			mcp.Description("菜单图标名称"),
			mcp.DefaultString("menu"),
		),
		mcp.WithBoolean("keepAlive",
			mcp.Description("是否缓存页面"),
		),
		mcp.WithBoolean("defaultMenu",
			mcp.Description("是否是基础路由"),
		),
		mcp.WithBoolean("closeTab",
			mcp.Description("是否自动关闭tab"),
		),
		mcp.WithString("activeName",
			mcp.Description("高亮菜单名称"),
		),
		mcp.WithString("parameters",
			mcp.Description("路由参数JSON字符串，格式：[{\"type\":\"params\",\"key\":\"id\",\"value\":\"1\"}]"),
		),
		mcp.WithString("menuBtn",
			mcp.Description("菜单按钮JSON字符串，格式：[{\"name\":\"add\",\"desc\":\"新增\"}]"),
		),
	)
}

// Handle 处理菜单创建请求
func (m *MenuCreator) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 解析请求参数
	args := request.GetArguments()

	// 必需参数
	path, ok := args["path"].(string)
	if !ok || path == "" {
		return nil, errors.New("path 参数是必需的")
	}

	name, ok := args["name"].(string)
	if !ok || name == "" {
		return nil, errors.New("name 参数是必需的")
	}

	component, ok := args["component"].(string)
	if !ok || component == "" {
		return nil, errors.New("component 参数是必需的")
	}

	title, ok := args["title"].(string)
	if !ok || title == "" {
		return nil, errors.New("title 参数是必需的")
	}

	// 可选参数
	parentId := uint(0)
	if val, ok := args["parentId"].(float64); ok {
		parentId = uint(val)
	}

	hidden := false
	if val, ok := args["hidden"].(bool); ok {
		hidden = val
	}

	sort := 1
	if val, ok := args["sort"].(float64); ok {
		sort = int(val)
	}

	icon := "menu"
	if val, ok := args["icon"].(string); ok && val != "" {
		icon = val
	}

	keepAlive := false
	if val, ok := args["keepAlive"].(bool); ok {
		keepAlive = val
	}

	defaultMenu := false
	if val, ok := args["defaultMenu"].(bool); ok {
		defaultMenu = val
	}

	closeTab := false
	if val, ok := args["closeTab"].(bool); ok {
		closeTab = val
	}

	activeName := ""
	if val, ok := args["activeName"].(string); ok {
		activeName = val
	}

	// 解析参数和按钮
	var parameters []system.SysBaseMenuParameter
	if parametersStr, ok := args["parameters"].(string); ok && parametersStr != "" {
		var paramReqs []MenuParameterRequest
		if err := json.Unmarshal([]byte(parametersStr), &paramReqs); err != nil {
			return nil, errors.Errorf("parameters 参数格式错误: %v", err)
		}
		for _, param := range paramReqs {
			parameters = append(parameters, system.SysBaseMenuParameter{
				Type:  param.Type,
				Key:   param.Key,
				Value: param.Value,
			})
		}
	}

	var menuBtn []system.SysBaseMenuBtn
	if menuBtnStr, ok := args["menuBtn"].(string); ok && menuBtnStr != "" {
		var btnReqs []MenuButtonRequest
		if err := json.Unmarshal([]byte(menuBtnStr), &btnReqs); err != nil {
			return nil, errors.Errorf("menuBtn 参数格式错误: %v", err)
		}
		for _, btn := range btnReqs {
			menuBtn = append(menuBtn, system.SysBaseMenuBtn{
				Name: btn.Name,
				Desc: btn.Desc,
			})
		}
	}

	// 构建菜单对象
	menu := system.SysBaseMenu{
		ParentId:  parentId,
		Path:      path,
		Name:      name,
		Hidden:    hidden,
		Component: component,
		Sort:      sort,
		Meta: system.Meta{
			Title:       title,
			Icon:        icon,
			KeepAlive:   keepAlive,
			DefaultMenu: defaultMenu,
			CloseTab:    closeTab,
			ActiveName:  activeName,
		},
		Parameters: parameters,
		MenuBtn:    menuBtn,
	}

	// 创建菜单
	menuService := service.ServiceGroupApp.SystemServiceGroup.MenuService
	err := menuService.AddBaseMenu(menu)
	if err != nil {
		return nil, errors.Errorf("创建菜单失败: %v", err)
	}

	// 获取创建的菜单ID
	var createdMenu system.SysBaseMenu
	err = global.GVA_DB.Where("name = ? AND path = ?", name, path).First(&createdMenu).Error
	if err != nil {
		logx.Error("获取创建的菜单ID失败", logx.Field("err", err))
	}

	// 构建响应
	response := &MenuCreateResponse{
		Success: true,
		Message: fmt.Sprintf("成功创建菜单 %s", title),
		MenuID:  createdMenu.ID,
		Name:    name,
		Path:    path,
	}

	resultJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return nil, errors.Errorf("序列化结果失败: %v", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("菜单创建结果：\n\n%s", string(resultJSON)),
			},
		},
	}, nil
}
