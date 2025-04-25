"""Histogram
"""

from collections import Counter
from pathlib import Path
import nltk

text = nltk.corpus.gutenberg.raw('austen-emma.txt')

counter = Counter(c for c in text if ord(c) < 128)
histogram = [counter.get(chr(i), 0) for i in range(128)]

print(counter)
print(histogram)

Path('histogram.txt').write_text('\n'.join(str(x) for x in histogram), encoding='utf8')
