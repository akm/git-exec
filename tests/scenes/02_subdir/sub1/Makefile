WORK_FILE=work.txt

.PHONY: add-one
add-one:
	@echo "One" >> $(WORK_FILE)

.PHONY: add-two
add-two:
	@echo "Two" >> $(WORK_FILE)

.PHONY: parent-%
parent-%:
	@$(MAKE) -C .. $*
