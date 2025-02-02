package events

import (
	"fmt"

	"github.com/mrpapercut/seca/models"
)

type FSDJumpEvent struct {
	GenericEvent
	StarSystem    string
	SystemAddress int64
	StarPos       []float64
	Body          string
	BodySystemID  int64 `json:"BodyID"`
	BodyType      string
	JumpDistance  float64 `json:"JumpDist"`
}

func (eh *EventHandler) handleEventFSDJump(rawEvent string) error {
	event, err := ParseEvent[FSDJumpEvent](rawEvent)
	if err != nil {
		return err
	}

	system := &models.System{
		Name:          event.StarSystem,
		SystemAddress: event.SystemAddress,
		StarPosX:      event.StarPos[0],
		StarPosY:      event.StarPos[1],
		StarPosZ:      event.StarPos[2],
		LastVisited:   event.Timestamp,
	}

	err = models.SaveSystem(system)
	if err != nil {
		return fmt.Errorf("error creating or updating system: %v", err)
	}

	retrievedSystem, err := models.GetSystemByAddress(event.SystemAddress)
	if err != nil {
		return fmt.Errorf("error retrieving saved system: %v", err)
	}

	body := &models.Body{
		Name:         event.Body,
		BodySystemID: event.BodySystemID,
		SystemID:     retrievedSystem.ID,
		System:       *retrievedSystem,
		BodyType:     event.BodyType,
	}

	err = models.SaveBody(body)
	if err != nil {
		return fmt.Errorf("error creating or updating body: %v", err)
	}

	status, err := models.GetStatus()
	if err != nil {
		return fmt.Errorf("error getting status: %v", err)
	}
	status.System = retrievedSystem.Name
	status.Body = body.Name

	status.TotalDistance += event.JumpDistance
	status.TotalJumps += 1

	err = models.UpdateStatus(status)
	if err != nil {
		return fmt.Errorf("error updating status: %v", err)
	}

	return nil
}
