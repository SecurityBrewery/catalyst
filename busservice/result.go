package busservice

import (
	"log"

	"github.com/SecurityBrewery/catalyst/bus"
	"github.com/SecurityBrewery/catalyst/generated/model"
)

func (h *busService) handleResult(resultMsg *bus.ResultMsg) {
	if resultMsg.Target != nil {
		ctx := busContext()
		switch {
		case resultMsg.Target.TaskOrigin != nil:
			if _, err := h.db.TaskComplete(
				ctx,
				resultMsg.Target.TaskOrigin.TicketId,
				resultMsg.Target.TaskOrigin.PlaybookId,
				resultMsg.Target.TaskOrigin.TaskId,
				resultMsg.Data,
			); err != nil {
				log.Println(err)
			}
		case resultMsg.Target.ArtifactOrigin != nil:
			enrichment := &model.EnrichmentForm{
				Data: resultMsg.Data,
				Name: resultMsg.Automation,
			}
			_, err := h.db.EnrichArtifact(ctx, resultMsg.Target.ArtifactOrigin.TicketId, resultMsg.Target.ArtifactOrigin.Artifact, enrichment)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
