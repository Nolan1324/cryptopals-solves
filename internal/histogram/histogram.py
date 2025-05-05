"""Histogram
"""

from collections import Counter
from pathlib import Path
from typing import Iterable
import nltk

def create_histogram(text: Iterable[str], output_file_name: str):
    counter = Counter(c for c in text if ord(c) < 128)
    histogram = [counter.get(chr(i), 0) for i in range(128)]

    print(counter)

    Path(output_file_name).write_text('\n'.join(str(x) for x in histogram), encoding='utf8')

corpus = nltk.corpus.gutenberg
# create_histogram(corpus.raw('austen-emma.txt'), 'histogram.txt')
create_histogram((sent[0][0] for sent in corpus.sents('austen-emma.txt')), 'histogram_first_chars.txt')
