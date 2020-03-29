package goq

// Join types.
const (
	JOIN_INNER = "INNER"
	JOIN_LEFT  = "LEFT OUTER"
	JOIN_RIGHT = "RIGHT OUTER"
	JOIN_FULL  = "FULL OUTER"
)

type JoinType string

type JoinDefiner interface {
	joinDef() *joinDef
}

type joinDef struct {
	table    TableLike
	on       PredExpr
	joinType JoinType
}

type JoinClause struct {
	joinType JoinType
	table    TableLike
}

func (jc *JoinClause) On(pred PredExpr) *JoinOn {
	return &JoinOn{&joinDef{jc.table, pred, jc.joinType}}
}

type JoinOn struct {
	def *joinDef
}

func (jo *JoinOn) joinDef() *joinDef {
	return jo.def
}

// JoinDef defines how to join a table.
// This uses 'INNER JOIN' by default.
type JoinDef struct {
	table    TableLike
	on       PredExpr
	joinType JoinType
}

func Join(table TableLike) *JoinDef {
	return &JoinDef{table: table, joinType: JOIN_INNER}
}

func (j *JoinDef) joinDef() *joinDef {
	return &joinDef{j.table, j.on, JOIN_INNER}
}

// On specifies a condition to join.
func (j *JoinDef) On(on PredExpr) *JoinDef {
	j.on = on
	return j
}

func (j *JoinDef) Left() *JoinDef {
	j.joinType = JOIN_LEFT
	return j
}

func (j *JoinDef) Right() *JoinDef {
	j.joinType = JOIN_RIGHT
	return j
}

func (j *JoinDef) Full() *JoinDef {
	j.joinType = JOIN_FULL
	return j
}
