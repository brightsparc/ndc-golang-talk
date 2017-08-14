PYTHON   = python
FASTTEXT = ../../fasttextgo/fasttext
MODEL    = model
PROBS    = 3

all: prep train test

help:
	@echo $(TAG)
	@echo "prep      - split speakers.jl into train.txt, test.txt"
	@echo "train     - fasttext supervised"
	@echo "test      - fasttext test"
	@echo "predict   - fasttext predict-prob"

prep:
	$(PYTHON) data_prep.py

train:
	time $(FASTTEXT) supervised -input train.txt -output $(MODEL) -wordNgrams 2 -dim 20 -epoch 1000 -lr 0.1

test:
	time $(FASTTEXT) test $(MODEL).bin test.txt $(PROBS)

predict:
	time $(FASTTEXT) predict-prob $(MODEL).bin test.txt $(PROBS) > predict.txt
	@echo
	tail -3 predict.txt
	@echo
	tail -3 test.txt | cut -b 1-150
