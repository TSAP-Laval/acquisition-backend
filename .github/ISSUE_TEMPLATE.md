## Please follow the guide below

- Issues submitted without this template format will likely be **ignored**.
- You will be asked some questions and requested to provide some information, please read them **carefully** and answer completely.
- Put an `x` into all the boxes [ ] relevant to your issue (like so [x]).
- Use the *Preview* tab to see how your issue will actually look like.

---

### Before submitting an issue, make sure you have:
- [ ] Read the [README](https://github.com/TSAP-Laval/acquisition-backend/blob/master/README.md)
- [ ] [Searched](https://github.com/TSAP-Laval/acquisition-backend/search?type=Issues) the bugtracker for similar issues including **closed** ones
- [ ] Reviewed the sample code in [_tests files](https://github.com/TSAP-Laval/acquisition-backend/tree/master/api)

### Purpose of your issue?
- [ ] Bug report (encountered problems/errors)
- [ ] Feature request (request for a new functionality)
- [ ] Question
- [ ] Other

---

### The following sections requests more details for particular types of issues, you can remove any section (the contents between the triple ---) not applicable to your issue.

---

### For a *bug report*, you **must** include *code* that will reproduce the error, and the *error log/traceback*.

Code:

```golang
 # Example code that will produce the error reported
 $ go test -v --race ./...
```

Error/Debug Log:

```golang
--- FAIL: TestBD (0.00s)
    equipe_test.go:20: Post http://localhost:3000/api/bd: dial tcp 127.0.0.1:3000: getsockopt: connection refused
panic: runtime error: invalid memory address or nil pointer dereference [recovered]
    panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x10 pc=0x49a731]

```

---

### Describe your issue

Explanation of your issue goes here. Please make sure the description is worded well enough to be understood with as much context and examples as possible.