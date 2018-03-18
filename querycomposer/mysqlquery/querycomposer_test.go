package mysqlquery

import "testing"

const (
	success = "\u2713"
	failed  = "\u2717"
)

func TestMysqlSelectQueryComposer_Compose(t *testing.T) {
	t.Log("Testing mysqlquerycomposer column. Input a struct")
	{
		selectComposer := NewMySqlSelectQueryComposer()

		in := struct {
			Id      int    `query:"id"`
			Name    string `query:"name"`
			Address string `query:"address"`
			Email   string `query:"email"`
			Phone   string `query:"phone"`
		}{}

		selectComposer.Columns(in)

		expectedColumnClause := `SELECT id, name, address, email, phone`

		if selectComposer.columnClause == expectedColumnClause {
			t.Logf("%s expected %s", success, selectComposer.columnClause)
		} else {
			t.Errorf("%s expected %s, but got %s", failed, expectedColumnClause, selectComposer.columnClause)
		}
	}

	t.Log("Testing mysqlquerycomposer column. Input a slice")
	{
		selectComposer := NewMySqlSelectQueryComposer()
		selectComposer.Columns([]string{
			"id", "product", "quantity", "so",
		})

		expectedColumnClause := "SELECT id, product, quantity, so"

		if selectComposer.columnClause == expectedColumnClause {
			t.Logf("%s expected %s", success, selectComposer.columnClause)
		} else {
			t.Errorf("%s expected %s got %s", failed, expectedColumnClause, selectComposer.columnClause)
		}
	}
}

func TestMysqlSelectQueryComposer(t *testing.T) {
	t.Log("Testing MySqlSelectQueryComposer")
	{
		selectComposer := NewMySqlSelectQueryComposer()

		data := struct {
			AgentId   int    `query:"agents.id"`
			AgentName string `query:"agents.name"`
			AgentNis  string `query:"agents.nis"`
			OrderId   int    `query:"orders.id"`
			OrderName string `query:"orders.name"`
		}{}

		query := selectComposer.Columns(data).
			PersistenceNames([]string{"agents", "orders"}).
			Where("orders.id = agents.id").
			OrderBy("orders.name desc, agents.nis").
			Paginate(2, 10).
			Compose()

		expectedQuery := `SELECT agents.id, agents.name, agents.nis, orders.id, orders.name FROM agents, orders WHERE orders.id = agents.id ORDER BY orders.name desc, agents.nis LIMIT 10 OFFSET 10`
		if query == expectedQuery {
			t.Logf("[%s] expected query [%s]", success, expectedQuery)
		} else {
			t.Errorf("[%s] expected query [%s], got [%s]", failed, expectedQuery, query)
		}
	}
}

func TestNewMySqlInsertQueryComposer(t *testing.T) {
	t.Log("Testing MySqlInsertQueryComposer")
	{
		insertComposer := NewMySqlInsertQueryComposer()

		data := struct {
			AgentId   int    `query:"id"`
			AgentName string `query:"name"`
			AgentNis  string `query:"nis"`
		}{}

		query := insertComposer.Columns(data).
			PersistenceNames([]string{"agents"}).
			Compose()
		expectedQuery := `INSERT INTO agents( id, name, nis ) VALUES ( ?,?,? )`
		if query == expectedQuery {
			t.Logf("[%s] expected query [%s]", success, expectedQuery)
		} else {
			t.Errorf("[%s] expected query [%s], but got [%s]", failed, expectedQuery, query)
		}
	}
}

func TestNewMySqlUpdateQueryComposer(t *testing.T) {
	t.Log("Testing MySqlUpdateQueryComposer")
	{
		updateComposer := NewMySqlUpdateQueryComposer()
		data := struct {
			AgentId   int    `query:"id,primary"`
			AgentName string `query:"name"`
			AgentNis  string `query:"nis"`
		}{}
		query := updateComposer.Columns(data).
			PersistenceNames([]string{"agents"}).
			Where("id=?").
			Compose()
		expectedQuery := `UPDATE agents SET name = ?, nis = ? WHERE id=?`

		if query == expectedQuery {
			t.Logf("[%s] expected query [%s]", success, expectedQuery)
		} else {
			t.Errorf("[%s] expected query [%s], got [%s]", failed, expectedQuery, query)
		}
	}
}
