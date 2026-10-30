.DEFAULT_GOAL := help

.PHONY: help call answer

# Support concise positional syntax: make call 1, make answer 1
ifeq ($(firstword $(MAKECMDGOALS)),call)
  CALL_ID := $(word 2,$(MAKECMDGOALS))
  $(eval $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS)):;@:)
endif

ifeq ($(firstword $(MAKECMDGOALS)),answer)
  ANSWER_ID := $(word 2,$(MAKECMDGOALS))
  $(eval $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS)):;@:)
endif

PHONE ?= +994501234567
STATUS ?= NEW
RESULT ?= ANSWERED

help:
	@echo "Available commands:"
	@echo "  make call [id]     - Start lead call workflow (default: 1)"
	@echo "  make answer [id]   - Simulate 3CX answered webhook (default: 1)"

call:
	@ID="$(or $(CALL_ID),$(ID),$(id),1)"; \
	curl -X POST http://localhost:8080/start-call \
	  -H "Content-Type: application/json" \
	  -d "{\"ID\": \"$$ID\", \"Phone\": \"$(PHONE)\", \"Status\": \"$(STATUS)\"}"
	@echo ""

answer:
	@ID="$(or $(ANSWER_ID),$(ID),$(id),1)"; \
	curl -X POST http://localhost:8080/3cx/webhook \
	  -H "Content-Type: application/json" \
	  -d "{\"leadId\": \"$$ID\", \"result\": \"$(RESULT)\"}"
	@echo ""
