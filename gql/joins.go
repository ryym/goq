package gql

const (
	JOIN_INNER = "INNER"
	JOIN_LEFT  = "LEFT OUTER"
	JOIN_RIGHT = "RIGHT OUTER"
	JOIN_FULL  = "FULL OUTER"
)

type JoinType string

type JoinClause struct {
	joinType JoinType
	table    TableLike
}

func (jc *JoinClause) On(pred PredExpr) JoinOn {
	return JoinOn{jc.table, pred, jc.joinType}
}

type JoinOn struct {
	Table TableLike
	On    PredExpr
	Type  JoinType
}
