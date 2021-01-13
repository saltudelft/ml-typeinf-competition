- To verify spec file:
  - Diff'ed new spec file with previous spec file
  - Verified that the two files are similar
  - Verified that new spec file has more lines

- To verify deduplicated dataset:
  - Retrieved # of files in original dataset (3114863 total, 636383 .py files)
  - Retrieved # of files in deduped dataset (2772150 total, 293670 .py files)
  - Retrieved # of lines in duplicate files list (354409 files)

- To verify file split:
  - Verified that number of lines in file split = number of code files in dataset
  - Train: 211 442, Test: 58 734, Valid: 23 494

```
Files found: 3114863
Code files found: 636383
```

```
Files found: 2772150
Code files found: 293670
```

```
*************************Loading all the tokenized Python source code files************************
************************************Preprocessing tokenized files************************************
Number of source code files: 510,289
Total number of tokens: 227,093,480
***********************Vectorize pre-processed source code files using TF-IDF**********************
**************************Building KNN index and finding nearest neighbors*************************
*******************************Finding exact and near duplicate files******************************
*********************Report duplication stats & saving detected duplicate files********************
Number of duplicated files: 400,245 (78.43%)
Number of detected clusters: 45,836
Avg. number of files per clones: 8.73
Median number of files per clones: 3.00
Duplication ratio: 69.45%
Finished duplicate files detection in 31.45 minutes.
```

```
duplicate_list_file = open("duplicate_files.txt", 'r')
duplicate_files = set([fname.strip() for fname in duplicate_list_file])

changed_dupes = [pathlib.Path(COPY_TARGET).joinpath(*pathlib.Path(fname).parts[1:]) for fname in duplicate_files]
exist_dupes = [fname for fname in changed_dupes if os.path.exists(fname)]
print(len(exist_dupes))
```

```
Splitting Python code files to train, test & validation sets
Number of files in train set: 211,442
Number of files in validation set: 23,494
Number of files in test set: 58,734
```