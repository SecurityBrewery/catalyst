package database

import (
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/SecurityBrewery/catalyst/caql"
	"github.com/SecurityBrewery/catalyst/dag"
	"github.com/SecurityBrewery/catalyst/generated/model"
)

func playbookGraph(playbook *model.Playbook) (*dag.Graph, error) {
	d := dag.NewGraph()

	var taskIDs []string
	for taskID := range playbook.Tasks {
		taskIDs = append(taskIDs, taskID)
	}
	sort.Strings(taskIDs)

	for _, taskID := range taskIDs {
		if err := d.AddNode(taskID); err != nil {
			return nil, errors.New("could not add node")
		}
	}
	for _, taskID := range taskIDs {
		task := playbook.Tasks[taskID]
		for next := range task.Next {
			if err := d.AddEdge(taskID, next); err != nil {
				return nil, errors.New("could not add edge")
			}
		}
	}
	return d, nil
}

func toTaskResponse(playbook *model.Playbook, taskID string, order int, graph *dag.Graph) (*model.TaskResponse, error) {
	task, ok := playbook.Tasks[taskID]
	if !ok {
		return nil, fmt.Errorf("task %s not found", taskID)
	}

	tr := &model.TaskResponse{
		Automation: task.Automation,
		Closed:     task.Closed,
		Created:    task.Created,
		Data:       task.Data,
		Done:       task.Done,
		Join:       task.Join,
		Payload:    task.Payload,
		Name:       task.Name,
		Next:       task.Next,
		Owner:      task.Owner,
		Schema:     task.Schema,
		Type:       task.Type,
		// Active:     active,
		// Order:      v.Order,
	}

	tr.Order = int64(order)

	taskActive, _ := active(playbook, taskID, graph, task)
	tr.Active = taskActive

	return tr, nil
}

func activePlaybook(playbook *model.Playbook, taskID string) (bool, error) {
	task, ok := playbook.Tasks[taskID]
	if !ok {
		return false, fmt.Errorf("playbook does not contain tasks %s", taskID)
	}

	d, err := playbookGraph(playbook)
	if err != nil {
		return false, err
	}

	return active(playbook, taskID, d, task)
}

func active(playbook *model.Playbook, taskID string, d *dag.Graph, task *model.Task) (bool, error) {
	if task.Done {
		return false, nil
	}

	parents := d.GetParents(taskID)

	if len(parents) == 0 {
		return true, nil // valid(&task)
	}

	if task.Join != nil && *task.Join {
		for _, parent := range parents {
			parentTask := playbook.Tasks[parent]
			if !parentTask.Done {
				return false, nil
			}
			requirement := parentTask.Next[taskID]

			b, err := evalRequirement(requirement, parentTask.Data)
			if err != nil {
				return false, err
			}

			if !b {
				return false, nil
			}
		}
		return true, nil
	}

	for _, parent := range parents {
		parentTask := playbook.Tasks[parent]
		if !parentTask.Done {
			// return false, nil
			continue
		}
		requirement := parentTask.Next[taskID]

		b, err := evalRequirement(requirement, parentTask.Data)
		if err != nil {
			continue
		}

		if b {
			return true, nil
		}
	}
	return false, nil
}

func evalRequirement(aql string, data interface{}) (bool, error) {
	if aql == "" {
		return true, nil
	}

	parser := caql.Parser{}
	tree, err := parser.Parse(aql)
	if err != nil {
		return false, err
	}

	var dataMap map[string]interface{}
	if data != nil {
		if dataMapX, ok := data.(map[string]interface{}); ok {
			dataMap = dataMapX
		} else {
			log.Println("wrong data type for task data")
		}
	}

	v, err := tree.Eval(dataMap)
	if err != nil {
		return false, err
	}

	if b, ok := v.(bool); ok {
		return b, nil
	}
	return false, err
}

/*
// "github.com/qri-io/jsonschema"
func valid(task *model.Task) (bool, error) {
	schema, err := json.Marshal(task.Schema)
	if err != nil {
		return false, err
	}

	rs := &jsonschema.Schema{}
	if err := json.Unmarshal(schema, rs); err != nil {
		return false, err
	}

	state := rs.Validate(context.Background(), task.Data)
	return len(*state.Errs) > 0, nil
}
*/
