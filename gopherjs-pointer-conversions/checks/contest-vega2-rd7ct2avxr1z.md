The verdict claims the tests require fixing "unrelated pre-existing GopherJS compiler bugs (& operator for promoted fields and closures, and reflect.Swapper)" that "were not specified". This is wrong on two points.

reflect.Swapper is not pre-existing. Swapping through (*[2]In2)(unsafe.Pointer(&in)) gives [{2} {1}] on unpatched GopherJS, native Go and the reference, but [{2} {2}] on this run. The bug comes from this run's own view proxies.

The closure and promoted-field failures are specified. The failing assertions compare converted pointers to one variable: "ptcIntPtr(&n) inside a function literal != ptcIntPtr(&n) outside of it", "conversions of &o.x and &o.ptcEmbedded.x are not equal", "unsafe.Pointer(&o.x) != unsafe.Pointer(&o.ptcEmbedded.x)". The description states: "pointers to the same variable must be equal after being converted to the same type, including unsafe.Pointer". These are stated requirements, not unrelated `&` behavior.
