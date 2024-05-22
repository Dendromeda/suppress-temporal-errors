# suppress-temporal-errors

A simple temporal logger middleware that lets you set some activity error types to be logged as info instead of error. This is intended to be used when you have certain activities using the built-in activity retry, for example when polling something. This way you will not pollute your error logs with expected stuff.

This function will create a new default temporal logger (as is in time of writing, the defaults are not exposed by sdk, the file is copied to this repo)
```
suppresserrors.NewLoggerWithSuppressedTypes("firstErrorType", "secondErrorType", "etc")
```

This function will use tour provided logger, with this as a middleware
```
suppresserrors.AddSuppressedErrorTypes(logger, "firstErrorType", "secondErrorType", "etc")
```

A super simple solution to a weird use case.
