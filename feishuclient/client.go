package feishuclient

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
)

type FeishuClient struct {
	*lark.Client
	pageSize int
}

func NewClient(appid, secret string) *FeishuClient {
	return &FeishuClient{
		lark.NewClient(appid, secret),
		10,
	}
}

func (fc *FeishuClient) SetPageSize(pageSize int) {
	fc.pageSize = pageSize
}

func (fc *FeishuClient) DepartmentChildrenQuery(id, idType, pageToken string, fetchChild bool) *larkcontact.ChildrenDepartmentResp {
	if idType == "" {
		idType = "open_department_id"
	}
	// 创建请求对象
	req := larkcontact.NewChildrenDepartmentReqBuilder().
		DepartmentId(id).
		UserIdType(`open_id`).
		DepartmentIdType(idType).
		FetchChild(fetchChild).
		PageToken(pageToken).
		PageSize(fc.pageSize).
		Build()

	// 发起请求
	resp, err := fc.Contact.V3.Department.Children(context.Background(), req)

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Printf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
		return nil
	}

	// 业务处理
	fmt.Println(larkcore.Prettify(resp))
	return resp
}

func (fc *FeishuClient) DepartmentChildrenQueryAll(id, idType string, fetchChild bool) []*larkcontact.Department {
	if idType == "" {
		idType = "open_department_id"
	}
	var ret []*larkcontact.Department
	pageToken := ""
	for {
		resp := fc.DepartmentChildrenQuery(id, idType, pageToken, fetchChild)
		if !resp.Success() {
			return ret
		}
		ret = append(ret, resp.Data.Items...)
		if *resp.Data.HasMore {
			pageToken = *resp.Data.PageToken
			continue
		}
		break
	}

	return ret
}
