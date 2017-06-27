# NDC Scraper

Python scraper to download the speak talks and tags for buliding a simple machine learning clasifier.

## Prerequisites

If you want to re-train the model or re-run the scrape you will require the following to be installed.

* python 2.6 or newer
* numpy & scipy
* sklearn
* pandas

## Installing

The machine learning model is built with `fastText`, and scrapy down the `Scrapy`.

### FastText

[fastText](https://github.com/facebookresearch/fastText) is a library for efficient learning of word representations and sentence classification.

You can d

```
$ git clone https://github.com/facebookresearch/fastText.git
$ cd fastText
$ make
```

### Scrapy

To run scrape locally you will need the [Scrapy](https://doc.scrapy.org/en/latest/intro/install.html) python library installed.

```
$ pip install Scrapy
```

You will then be able to run the scrape with the following command.

```
$ scrapy crawl speakers -o speakers.jl
```

This will create or append JSON lines for each speaker to the file `speakers.jl`

## Building Model

Explain how to run the automated tests for this system

### Data Prep

Data prep will create two text files `tain.txt` and `test.txt` in a format suitable for `fastText`.

```
$ python data_prep.py
```

### Train and Test Model

Train supervised model from the taining data annotated with tags as labels.

```
$ fasttext supervised -input train.txt -output model -wordNgrams 2 -lr 0.1 -epoch 1000
Read 0M words
Number of words:  2763
Number of labels: 29
Progress: 100.0%  words/sec/thread: 1594669  lr: 0.000000  loss: 2.017118  eta: 0h0m
```

Evaluate it by computing the precision and recall at k (P@k and R@k) on a test set

```
$ fasttext test model.bin test.txt 3
N	25
P@3	0.2
R@3	0.263
Number of examples: 25
```

Get the predict labels back test set, and inspect some records

```
$ fasttext predict model.bin test.txt 3 | head -n 3
__label__ux __label__machine __label__cloud
__label__devops __label__continuous __label__cloud
__label__ux __label__architecture __label__cloud
```

## Authors

* **Julian Bright**

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Thanks for SEEKers who helped pull this together for NDC.
