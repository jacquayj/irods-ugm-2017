TestBasicExample {
    msibasic_example("keytest,valuetest", *outKVP);
    msiPrintKeyValPair("stderr", *outKVP)
}

INPUT null
OUTPUT ruleExecOut