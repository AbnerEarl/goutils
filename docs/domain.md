
## Usage

This is a simple example that demonstrates how to use the package with the default options and the default Public Suffix list packaged with the library.

```go
package main

import (
    "fmt"

    "github.com/AbnerEarl/goutils/domain"
)

func main() {
    // Extract the domain from a string
    // using the default list
    fmt.Println(domain.Domain("example.com"))             // example.com
    fmt.Println(domain.Domain("www.example.com"))         // example.com
    fmt.Println(domain.Domain("example.co.uk"))           // example.co.uk
    fmt.Println(domain.Domain("www.example.co.uk"))       // example.co.uk

    // Parse the domain from a string
    // using the default list
    fmt.Println(domain.Parse("example.com"))             // &DomainName{"com", "example", ""}
    fmt.Println(domain.Parse("www.example.com"))         // &DomainName{"com", "example", "www"}
    fmt.Println(domain.Parse("example.co.uk"))           // &DomainName{"co.uk", "example", ""}
    fmt.Println(domain.Parse("www.example.co.uk"))       // &DomainName{"co.uk", "example", "www"}
}
```

#### Ignoring Private Domains

The PSL is composed by two list of suffixes: IANA suffixes, and Private Domains.

Private domains are submitted by private organizations. By default, private domains are not ignored.
Sometimes, you want to ignore these domains and only query against the IANA suffixes. You have two options:

1. Ignore the domains at runtime
2. Create a custom list without the private domains

In the first case, the private domains are ignored at runtime: they will still be included in the lists but the lookup will skip them when found.

```go
domain.DomainFromListWithOptions(domain.DefaultList(), "google.blogspot.com", nil)
// google.blogspot.com

domain.DomainFromListWithOptions(domain.DefaultList(), "google.blogspot.com", &domain.FindOptions{IgnorePrivate: true})
// blogspot.com

// Note that the DefaultFindOptions includes the private domains by default
domain.DomainFromListWithOptions(domain.DefaultList(), "google.blogspot.com", domain.DefaultFindOptions)
// google.blogspot.com
```

This solution is easy, but slower. If you find yourself ignoring the private domains in all cases (or in most cases), you may want to create a custom list without the private domains.

```go
list := NewListFromFile("path/to/list.txt", &domain.ParserOption{PrivateDomains: false})
domain.DomainFromListWithOptions(list, "google.blogspot.com", nil)
// blogspot.com
```

## IDN domains, A-labels and U-labels

[A-label and U-label](https://tools.ietf.org/html/rfc5890#section-2.3.2.1) are two different ways to represent IDN domain names. These two encodings are also known as ASCII (A-label) or Pynucode vs Unicode (U-label). Conversions between U-labels and A-labels are performed according to the ["Punycode" specification](https://tools.ietf.org/html/rfc3492), adding or removing the ACE prefix as needed.

IDNA-aware applications generally use the A-label form for storing and manipulating data, whereas the U-labels can appear in presentation and user interface forms.

Although the PSL list has been traditionally U-label encoded, this library follows the common industry standards and stores the rules in their A-label form. Therefore, unless explicitly mentioned, any method call, comparison or internal representation is expected to be ASCII-compatible encoded (ACE).

Passing Unicode names to the library may either result in error or unexpected behaviors.

If you are interested in the details of this decision, you can read the full discussion [here](https://github.com/AbnerEarl/goutils/issues/31).


## Differences with `golang.org/x/net/publicsuffix`

The [`golang.org/x/net/publicsuffix`](https://godoc.org/golang.org/x/net/publicsuffix) is a package part of the Golang `x/net` package, that provides a public suffix list implementation.

The main difference is that the `x/net` package is optimized for speed, but it's less flexible. The list is compiled and embedded into the package itself. However, this is also the main downside.
The [list is not frequently refreshed](https://github.com/letsencrypt/boulder/issues/1374#issuecomment-182429297), hence the results may be inaccurate, in particular if you heavily rely on the private domain section of the list. Changes in the IANA section are less frequent, whereas changes in the Private Domains section happens weekly.

This package provides the following extra features:

- Ability to load an arbitrary list at runtime (e.g. you can feed your own list, or create multiple lists)
- Ability to create multiple lists
- Ability to parse a domain using a previously defined list
- Ability to add custom rules to an existing list, or merge/load rules from other lists (provided as file or string)
- Advanced access to the list rules
- Ability to ignore private domains at runtime, or when the list is parsed

This package also aims for 100% compatibility with the `x/net` package. A special adapter is provided as a drop-in replacement. Simply change the include statement from

```go
import (
    "golang.org/x/net/publicsuffix"
)
```

to

```go
import (
    "github.com/AbnerEarl/goutils/net/domain"
)
```

The `github.com/AbnerEarl/goutils/net/domain` package defines the same methods defined in `golang.org/x/net/publicsuffix`, but these methods are implemented using the `github.com/AbnerEarl/goutils/domain` package.

Note that the adapter doesn't offer the flexibility of `github.com/AbnerEarl/goutils/domain`, such as the ability to use multiple lists or disable private domains at runtime.


## `cookiejar.domainList` interface

This package implements the [`cookiejar.domainList` interface](https://godoc.org/net/http/cookiejar#domainList). It means it can be used as a value for the `domainList` option when creating a `net/http/cookiejar`.

```go
import (
    "net/http/cookiejar"
    "github.com/AbnerEarl/goutils/domain"
)

deliciousJar := cookiejar.New(&cookiejar.Options{domainList: domain.CookieJarList})
```
