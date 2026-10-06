defmodule SampleTest do
 use ExUnit.Case
 test "actual umbrella value", do: assert(Sample.value() == 42)
end
