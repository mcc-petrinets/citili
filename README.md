# Citili
Generation of CTL formulas for the model checking contest

## Requirements

The tool is based on the [Pinimili](https://github.com/loig/pinimili) parser for PNML. It also uses a [simple model checker](https://github.com/mcc-petrinets/formulas/tree/master/smc) that was developped for earlier generators for filtering formulas (only needed at execution, not for compiling the tool). We plan to no longer rely on this model checker in the futur.
