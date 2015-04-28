package sql

import (
	"github.com/andrew-d/docstore/query/parse"
)

func NodesToSQL(root parse.Node, defaultField string) (sql string, args []interface{}) {
	var traverse func(node parse.Node) (string, []interface{})

	traverse = func(node parse.Node) (string, []interface{}) {
		switch v := node.(type) {
		case *parse.TextNode:
			fname := v.Field
			if fname == "" {
				fname = defaultField
			}

			return fname + " = ?", []interface{}{v.Text}

		case *parse.OrNode:
			leftSql, leftArgs := traverse(v.Left)
			rightSql, rightArgs := traverse(v.Right)

			sql := "(" + leftSql + ") OR (" + rightSql + ")"
			args := make([]interface{}, len(leftArgs)+len(rightArgs))
			copy(args[0:], leftArgs)
			copy(args[len(leftArgs):], rightArgs)

			return sql, args

		case *parse.AndNode:
			leftSql, leftArgs := traverse(v.Left)
			rightSql, rightArgs := traverse(v.Right)

			sql := "(" + leftSql + ") AND (" + rightSql + ")"
			args := make([]interface{}, len(leftArgs)+len(rightArgs))
			copy(args[0:], leftArgs)
			copy(args[len(leftArgs):], rightArgs)

			return sql, args

		case *parse.NotNode:
			sql, args := traverse(v.Node)

			return "NOT (" + sql + ")", args

		default:
			panic("unknown node type")
		}
	}

	return traverse(root)
}
