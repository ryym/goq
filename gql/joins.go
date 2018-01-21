package gql

const (
	JOIN_INNER = "INNER"
	JOIN_LEFT  = "LEFT OUTER"
	JOIN_RIGHT = "RIGHT OUTER"
	JOIN_FULL  = "FULL OUTER"
)

type JoinType string

type JoinDefiner interface {
	joinDef() *JoinDef
}

type JoinDef struct {
	Table TableLike
	On    PredExpr
	Type  JoinType
}

type JoinClause struct {
	joinType JoinType
	table    TableLike
}

func (jc *JoinClause) On(pred PredExpr) *JoinOn {
	return &JoinOn{&JoinDef{jc.table, pred, jc.joinType}}
}

type JoinOn struct {
	def *JoinDef
}

func (jo *JoinOn) joinDef() *JoinDef {
	return jo.def
}

type Join struct {
	Table TableLike
	On    PredExpr
}

func (j *Join) joinDef() *JoinDef {
	return &JoinDef{j.Table, j.On, JOIN_INNER}
}

func (j *Join) Left() *JoinDef {
	return &JoinDef{j.Table, j.On, JOIN_LEFT}
}

func (j *Join) Right() *JoinDef {
	return &JoinDef{j.Table, j.On, JOIN_RIGHT}
}

func (j *Join) Full() *JoinDef {
	return &JoinDef{j.Table, j.On, JOIN_FULL}
}
