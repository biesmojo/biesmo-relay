package engine

import (
	"context"

	"relay/api/models"
)

func EvaluateRules(pool *pgxpool.Pool, event models.Event) {
	rows, err := pool.Query(context.Background(),
		"SELECT id, name, match_conditions, action_type, action_payload FROM rules WHERE event_type = $1 AND is_active = true",
		event.Type)
	if err != nil {
		log.Printf("Failed to load rules: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var rule models.Rule
		err := rows.Scan(&rule.ID, &rule.Name, &rule.MatchConditions, &rule.ActionType, &rule.ActionPayload)
		if err != nil {
			continue
		}

		if matches(event.Payload, rule.MatchConditions) {
			go ExecuteAction(pool, rule, event)
			logRuleAudit(pool, rule.ID, event.ID, rule.ActionType, "matched", nil)
		} else {
			logRuleAudit(pool, rule.ID, event.ID, rule.ActionType, "no_match", nil)
		}
	}
}

func matches(payload, conditions json.RawMessage) bool {
	var p map[string]interface{}
	json.Unmarshal(payload, &p)
	var c map[string]interface{}
	json.Unmarshal(conditions, &c)

	for k, v := range c {
		pv, ok := p[k]
		if !ok {
			return false
		}
		switch cv := v.(type) {
		case []interface{}:
			found := false
			for _, item := range cv {
				if item == pv {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		default:
			if v != pv {
				return false
			}
		}
	}
	return true
}

func logRuleAudit(pool *pgxpool.Pool, ruleID, eventID int, actionType string, outcome string, errMsg *string) {
	_, _ = pool.Exec(context.Background(),
		"INSERT INTO rule_audit_log (rule_id, event_id, action_type, outcome, error_msg) VALUES ($1, $2, $3, $4, $5)",
		ruleID, eventID, actionType, outcome, errMsg)
}
