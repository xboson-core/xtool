# xtool

`xtool` is a file packaging and incremental update tool.

It can build a complete ZIP package from a directory, generate an incremental ZIP package by comparing a directory against an existing package metadata file, and record removed files for later processing.

## Features

* Build a complete ZIP package.
* Generate incremental ZIP packages.
* Track file size and modification time through a CSV metadata file.
* Record removed files.
* Exclude files using a rule file.
* Show processing progress.
* Support single-line progress output.
* SQL diff mode is reserved for future use.

## Usage

```text
xtool <mode> [options]
```

Exactly one mode must be specified.

### Modes

| Mode    | Description                       |
| ------- | --------------------------------- |
| `-pack` | Build a complete ZIP package      |
| `-diff` | Build an incremental ZIP package  |
| `-sqld` | SQL diff mode, currently reserved |

---

## Pack Mode

Use `-pack` to create a complete ZIP package from an input directory.

```bash
xtool -pack -i=/data/ui/web -o=web.zip
```

This will:

1. Scan all files under `/data/ui/web`.
2. Add all files to `web.zip`.
3. Generate the metadata file `web.zip.csv` by default.

The metadata contains:

```text
file path, file size, modification time
```

For example:

```csv
/data/ui/web/index.html,1234,1757654321
/data/ui/web/js/app.js,45678,1757654300
```

### Custom Metadata File

Use `-meta` to specify the metadata file:

```bash
xtool -pack \
  -i=/data/ui/web \
  -o=web.zip \
  -meta=web-meta.csv
```

If `-meta` is not specified, the default is:

```text
<output>.csv
```

For example:

```text
web.zip
web.zip.csv
```

---

## Diff Mode

Use `-diff` to generate an incremental package.

```bash
xtool -diff \
  -i=/data/ui/web \
  -b=web.zip \
  -o=update.zip
```

The `-b` option specifies the existing base ZIP package.

The corresponding metadata file is used to determine which files have changed.

By default:

```text
meta = <base>.csv
```

Therefore:

```bash
xtool -diff -i=/data/ui/web -b=web.zip -o=update.zip
```

implicitly uses:

```text
web.zip.csv
```

### How Diff Works

For each file in the input directory, `xtool` compares its current metadata with the metadata stored in the base package.

A file is included in the incremental package when its metadata has changed.

The metadata includes:

* File path
* File size
* Unix modification time

Unchanged files are not included in the incremental ZIP.

The original metadata file is **not modified** by diff mode.

### Removed Files

Files that exist in the base metadata but no longer exist in the input directory can be recorded using `-del`.

```bash
xtool -diff \
  -i=/data/ui/web \
  -b=web.zip \
  -o=update.zip \
  -del=removed.txt
```

The resulting file contains the paths of files removed from the input directory.

---

## Excluding Files

Use `--exclude-from` to specify a file containing exclusion rules.

```bash
xtool -pack \
  -i=/data/ui/web \
  -o=web.zip \
  --exclude-from=exclude.txt
```

The format follows the same general pattern as the `tar --exclude-from` option.

For example:

```text
*.log
*.tmp
node_modules
.git
```

Excluded files are not processed or added to the package.

---

## Progress Output

Use `-v` to display processing progress:

```bash
xtool -pack \
  -i=/data/ui/web \
  -o=web.zip \
  -v
```

Use `-l` to display progress on a single line:

```bash
xtool -pack \
  -i=/data/ui/web \
  -o=web.zip \
  -l
```

`-l` enables single-line progress output.

---

## Command-Line Options

| Option           | Description                      |
| ---------------- | -------------------------------- |
| `-pack`          | Build a complete ZIP package     |
| `-diff`          | Build an incremental ZIP package |
| `-sqld`          | SQL diff mode                    |
| `-i`             | Input directory                  |
| `-o`             | Output ZIP file                  |
| `-b`             | Base ZIP file used by diff mode  |
| `-meta`          | Metadata CSV file                |
| `--exclude-from` | File containing exclusion rules  |
| `-del`           | Output file for removed files    |
| `-v`             | Show processing progress         |
| `-l`             | Show progress on a single line   |

Exactly one of `-pack`, `-diff`, or `-sqld` must be specified.

---

## Metadata

The metadata CSV is used to determine whether files have changed.

Each record contains three fields:

```text
path,size,unix
```

Example:

```csv
/data/ui/web/index.html,1234,1757654321
/data/ui/web/css/app.css,5678,1757654320
/data/ui/web/js/app.js,10240,1757654310
```

The metadata is intentionally simple so it can be generated, inspected, and processed by other tools.

No file hash is currently required for change detection.

---

## Typical Workflow

### 1. Build the initial package

```bash
xtool -pack \
  -i=/data/ui/web \
  -o=web.zip
```

This produces:

```text
web.zip
web.zip.csv
```

### 2. Modify the source directory

Files can be added, modified, or removed from:

```text
/data/ui/web
```

### 3. Build an incremental update

```bash
xtool -diff \
  -i=/data/ui/web \
  -b=web.zip \
  -o=update.zip \
  -del=removed.txt
```

The result contains only the files that need to be updated, while `removed.txt` contains files that no longer exist.


## SQL Diff

The `--sqld` mode generates a SQL migration script by comparing a **base SQL file** with an **updated SQL file**.

In this mode:

* `-b` specifies the base SQL file.
* `-i` specifies the updated SQL file.
* `-o` specifies the output SQL file.
* `-df` specifies an optional SQL Diff configuration file.

The generated SQL contains the changes required to migrate the base database to the updated database.

### Example

```bash
./xtool --sqld \
    -b base.sql \
    -i updated.sql \
    -o migration.sql \
    -df sqlconf.yaml
```

This compares the SQL dump specified by `-b` with the updated SQL dump specified by `-i`, and writes the generated migration SQL to `-o`.

### SQL Diff Configuration

The `-df` option can be used to exclude specific tables or schemas from the comparison.

This is useful when the SQL dumps contain tables that are generated dynamically, contain runtime logs, or belong to system databases that should not be included in the migration.

The configuration file uses YAML format:

```yaml
skip:
  tables:
    - a
    - b
  schemas:
    - mysql
    - db1
    - db2
```

#### `skip.tables`

`skip.tables` specifies table names that should be excluded from SQL Diff.

These tables are ignored during comparison, so changes to their structure or data will not be included in the generated migration SQL.

#### `skip.schemas`

`skip.schemas` specifies database schemas that should be excluded from SQL Diff.

All tables belonging to these schemas are ignored during comparison.

### mysqldump Compatibility

SQL Diff is designed specifically for SQL files exported by the **`mysqldump`** application.

It is **not a general-purpose SQL diff tool** and should not be used to compare arbitrary SQL scripts.

The tool automatically ignores dump-related statements and metadata that are not relevant to database migration, such as MySQL session configuration statements generated by `mysqldump`.

For example, statements such as:

```sql
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET FOREIGN_KEY_CHECKS=0 */;
```

are not treated as database changes and are not included in the generated migration SQL.

Comments and other `mysqldump`-specific metadata are also ignored where appropriate.

The purpose of this mode is to generate a **minimal migration script** containing only the database changes required to transform the base database into the updated database.

The typical workflow is:

```text
Base mysqldump
      │
      ▼
   SQL Diff
      ▲
      │
Updated mysqldump
      │
      ▼
Migration SQL
```

For example:

```bash
./xtool --sqld \
    -b base.sql \
    -i updated.sql \
    -o migration.sql \
    -df sqlconf.yaml
```

The resulting `migration.sql` can then be used to upgrade a database initialized from `base.sql` to the state represented by `updated.sql`.


## Build

Build the executable with Go:

```bash
go build -o xtool
```

On Windows:

```bash
go build -o xtool.exe
```

## License

See the project license for details.
