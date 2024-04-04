## Diff URLs

Remove duplicate URLs by retaining only the unique combinations of hostname, path, and parameter names.

## Install

```bash
go install github.com/j3ssie/durl@latest
```


## Usage

```bash
cat wayback_urls.txt | durl | tee differ_urls.txt

# with extra regex
cat wayback_urls.txt | durl -e 'your-regex-here' | tee differ_urls.txt
```

## Covered cases

The following examples illustrate the criteria used to ensure each URL is considered unique and listed only once:

- URLs with the same hostname, path, and parameter names

```
http://sample.example.com/product.aspx?productID=123&type=customer
http://sample.example.com/product.aspx?productID=456&type=admin
```

- Paths indicating static content like blog, news or calender.

```
https://www.example.com/cn/news/all-news/public-1.html
https://www.sample.com/de/about/business/countrysites.htm
https://www.sample.com/de/about/business/very-long-string-here-that-exceed-100-char.htm
https://www.sample.com/de/blog/2022/01/02/blog-title.htm
```

- URLs with numeric variations

```
https://www.example.com/data/0001.html
https://www.example.com/data/0002.html
```

- Static file will be ignore like `http://example.com.com/cdn-cgi/style.css`
