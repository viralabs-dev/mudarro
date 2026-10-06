using System;
internal static class Program {
 static int Main(string[] args) {
  int actual=21*2;
  if(args.Length==1 && args[0]=="--controlled-failure") { Console.WriteLine("DOTNET_CONTROLLED_FAILURE_EXPECTED_EXIT_7");return 7; }
  if(args.Length==1 && args[0]=="--self-test") {if(actual!=42){Console.Error.WriteLine("SELF_TEST_FAIL");return 9;}Console.WriteLine("DOTNET_REAL_STDLIB_SELF_TEST_PASS: 21 * 2 = 42");return 0;}
  Console.WriteLine("DOTNET_REAL_CONSOLE_OUTPUT: "+actual);return 0;
 }
}
