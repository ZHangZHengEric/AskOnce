package utils

import "math"

func SlicePage(page, pageSize, nums int) (sliceStart, sliceEnd int) {
	if page <= 0 {
		page = 1
	}
	if pageSize < 0 {
		pageSize = 20 //设置一页默认显示的记录数
	}
	if pageSize > nums {
		return 0, nums
	}
	// 总页数
	pageCount := int(math.Ceil(float64(nums) / float64(pageSize)))
	if page > pageCount {
		return 0, 0
	}
	sliceStart = (page - 1) * pageSize
	sliceEnd = sliceStart + pageSize

	if sliceEnd > nums {
		sliceEnd = nums
	}
	return sliceStart, sliceEnd
}
