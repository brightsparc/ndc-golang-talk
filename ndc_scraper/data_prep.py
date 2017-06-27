import re
import pandas as pd
from sklearn.model_selection import train_test_split

def clean(t, replace=' '):
    return re.sub('\W+', ' ', t).lower()

# Load the lines and set label and contents
df = pd.read_json('speakers.jl', lines=True)
df['talk_labels'] = df.talk.apply(lambda talk: ' '.join(["__label__%s" % clean(t, '') for t in talk["tags"]])).values
df['talk_contents'] = df.talk.apply(lambda talk: ' '.join([clean(talk["preamble"]), clean(talk["body"])])).values

# Train test split
train, test = train_test_split(df[['talk_labels', 'talk_contents']].dropna(), test_size = 0.2, random_state=42)

# Write back labels and conents
train.to_csv('train.txt', sep='\t', header=False, index=False)
test.to_csv('test.txt', sep='\t', header=False, index=False)
