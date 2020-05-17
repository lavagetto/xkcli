# xkcli - set of tools for interacting with XKCD from the CLI

Have you ever dreamt of being able to search for a relevant XKCD strip in the middle of a conversation on IRC?

Enter XKCLI, the cli tool no one felt was needed, but you know, the author wants a toy project.

## Installation
Clone this repository and run
```
xkcli $ go mod vendor
xkcli $ go build .
```

## Usage
XKCLI works by creating a local index of strips, and then querying it.

First of all, you should refresh your database of local strips:
```
~ $ xkcli refresh -c 3
```
Here we're choosing a concurrency of 3 parallel downloads with the `-c` flag. By default, this command will create your index at `~/.xkcli.db`. This value can be changed by providing your configuration (see below).

Once your database is updated, you can search the strips as follows:
```
$ xkcli search "Star Trek"
Your search results:
0 - (1.95) XKCD 1167 (2013-01-30): Star Trek into Darkness
	strip: https://imgs.xkcd.com/comics/star_trek_into_darkness.png
1 - (1.13) XKCD 2041 (2018-09-03): Frontiers
	strip: https://imgs.xkcd.com/comics/frontiers.png
2 - (1.07) XKCD 902 (2011-05-23): Darmok and Jalad
	strip: https://imgs.xkcd.com/comics/darmok_and_jalad.png
3 - (0.99) XKCD 465 (2008-08-20): Quantum Teleportation
	strip: https://imgs.xkcd.com/comics/quantum_teleportation.png
4 - (0.87) XKCD 1429 (2014-10-03): Data
	strip: https://imgs.xkcd.com/comics/data.png
5 - (0.71) XKCD 1313 (2014-01-06): Regex Golf
	strip: https://imgs.xkcd.com/comics/regex_golf.png
We also found 4 results below the threshold (0.50)
```

If you're interested in just the link to the most relevant strip (for instance for use in an IRC client or similar IM system that allows running commands), you can use the `--lucky|-l` flag:

```
$ xkcli search -l "bobby tables"
https://xkcd.com/327
```

## Configuration
You can override the default behaviour of xkcli by editing `~/.xkcli.yaml`.

The currently implemented configurations are listed below:

| Name     | Description                     |     Default     |  Type  |
| :------- | :------------------------------ | :-------------: | :----: |
| minScore | Minimum score of search results |       0.5       | float  |
| dbPath   | Full path of the db directory   | $HOME/.xkcli.db | string |

## FAQ

### Do you have any relationship with XKCD?
None whatsoever. I just find the strips fun. In fact, I don't think this violates anyone's intellectual property, and was done just for fun as a weekend project to exercise myself with golang (and it shows!).

### Couldn't one just search the strip on Google? Isn't that much better?
Yes, in fact you should.

### Can you add feature X?
Maybe. This project is a toy written over weekends, new features might be implemented by the original authors if they find it interesting. You're welcome to add it yourself though. Or, you know, use your favourite search engine instead.

### Can you fix bug X?
Maybe - see the previous answer for a rationale. Specifically, bugs regarding the quality of search are unlikely to be fixed quickly.

### How do I run xkcli as a service on kubernetes?
LOL.