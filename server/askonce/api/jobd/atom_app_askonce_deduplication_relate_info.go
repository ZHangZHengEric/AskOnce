package jobd

func (entity *JobdApi) DeduplicationRelateInfo(input *GenerateRelateInfoRes) (res *GenerateRelateInfoRes, err error) {
	return doTaskProcess[*GenerateRelateInfoRes, *GenerateRelateInfoRes](entity, "atom_app_askonce_deduplication_relate_info", input, 100000)
}
