package dto

/**
 * @author Rancho
 * @date 2022/1/6
 */

type Pager struct {
    Page      int `json:"page"`
    PageSize  int `json:"page_size"`
    TotalRows int `json:"total_rows"`
}
