package sql

import (
	"fmt"
	"strings"
)

// SELECT * FROM users;
type SqlParser struct {
}

func (parser *SqlParser) parse(sql string) {
	var keywords = strings.Split(sql, " ")
	fmt.Println(keywords)
}
