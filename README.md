## Diff URLs

Remove duplicate URLs by retaining only the unique combinations of hostname, path, and parameter names.

## Install

```bash
go install github.com/j3ssie/durl@latest
```

## Usage

```bash
# basic usage
cat wayback_urls.txt | durl | tee differ-urls.txt

# with extra regex
cat wayback_urls.txt | durl -e 'your-regex-here' | tee differ-urls.txt

# only get the scope domain
cat spider-urls.txt | durl -t 'target.com' | tee in-scope-url.txt

# parse JSONL data
cat large-jsonl-data.txt | durl -t 'target.com' -f url | tee in-scope-jsonl-data.txt
```

## Covered cases

The following examples illustrate the criteria used to ensure each URL is considered unique and listed only once:

1. URLs with the same hostname, path, and parameter names

```
http://sample.example.com/product.aspx?productID=123&type=customer
http://sample.example.com/product.aspx?productID=456&type=admin
```

2. Paths indicating static content like blog, news or calender.

```
https://www.example.com/cn/news/all-news/public-1.html
https://www.sample.com/de/about/business/countrysites.htm
https://www.sample.com/de/about/business/very-long-string-here-that-exceed-100-char.htm
https://www.sample.com/de/blog/2022/01/02/blog-title.htm
```

3. URLs with numeric variations

```
https://www.example.com/data/0001.html
https://www.example.com/data/0002.html
```

4. Static file will be ignore like `http://example.com.com/cdn-cgi/style.css`

5. Select a url JSON field from the input then filtering with all of the cases above.