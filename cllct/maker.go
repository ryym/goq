package cllct

func NewMaker() *CollectorMaker {
	return &CollectorMaker{}
}

type CollectorMaker struct{}

func (cm *CollectorMaker) ToRowMapSlice(slice *[]map[string]interface{}) *RowMapSliceCollector {
	return &RowMapSliceCollector{slice: slice}
}
