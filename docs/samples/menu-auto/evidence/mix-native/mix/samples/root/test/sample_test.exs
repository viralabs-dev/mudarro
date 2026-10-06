defmodule SampleTest do
 use ExUnit.Case
 test "actual value", do: assert(Sample.value() == 42)
end
