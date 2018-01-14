package cllct

func NewMaker() *CollectorMaker {
	return &CollectorMaker{}
}

type CollectorMaker struct{}

func (cm *CollectorMaker) ToRowMapSlice(slice *[]map[string]interface{}) *RowMapSliceCollector {
	return &RowMapSliceCollector{slice: slice}
}

func (cm *CollectorMaker) ToRowMap(mp *map[string]interface{}) *RowMapCollector {
	return &RowMapCollector{mp: mp}
}
