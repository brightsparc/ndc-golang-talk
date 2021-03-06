import re
import pandas as pd
from sklearn.model_selection import train_test_split

def clean(t, replace=' '):
    if t:
        t = t.lower().replace('.net', 'dotnet').replace('c++', 'cpp') # special for nonword characters
        return re.sub('\s+', re.sub('\W+', replace, t), ' ')
    return ''

# Load data, and set talk properties
df = pd.read_json('speakers.jl', lines=True)
# Load the lines and set label and contents
df['talk_labels'] = df.talk.apply(lambda talk: ' '.join(["__label__%s" % clean(t, '_') for t in talk["tags"]]))
df['talk_contents'] = df.talk.apply(lambda talk: ' '.join([clean(talk["preamble"]), clean(talk["body"])]))

# Make sure we have labels and some contents with some length
df['talk_labels_len'] = df['talk_labels'].apply(lambda labels: len(labels))
df['talk_contents_len'] = df['talk_contents'].apply(lambda contents: len(contents))
df = df[(df['talk_labels_len']>5)&(df['talk_contents_len']>20)]

# Train test split
train, test = train_test_split(df[['talk_labels', 'talk_contents']], stratify=df['conference'], test_size = 0.2, random_state=42)

# Write back labels and contents (use tab so text is not escaped)
train.to_csv('train.txt', sep='\t', header=False, index=False)
test.to_csv('test.txt', sep='\t', header=False, index=False)
