
# Overview

This document specifies requirements of the DUCK compliance checker, with user stories, and provides a glossary of technical terms relevant for the compliance checker.  The document concludes with a list of next tasks.

# Relevant Objectives from the Contract

5. (P2) To provide compliance checks regarding legal regulations on data privacy (parametric)

6. (P2) To provide recommendations and extensions in case a set of data use declarations is not compliant to a given set of regulations.

7. (P2) All input functions and user interfaces are to be done in a localizable manner.

8. (P2) The tool’s design and architecture needs to allow for hosting as a service available to public.

9. (P2) The source code of DUCK will be published under an Open Source License.


# User Roles

As per Jim's email from 2016-05-03.

**Author**. Creates and manages data use documents

**Modeler**. Creates and manages rule bases

**Administrator**. Deploys and maintains the DUCK system

**Developer**. Customizes and extends the DUCK system

# Glossary

**Argument**. A (premises, conclusion) pair, where the premises are a set of statements and the conclusion is a statement.  The argument asserts that the premises provide reasons for accepting the conclusion. Arguments can be attacked and defeated by counterarguments of various kinds.

**Argument evaluation**. A method for determining the acceptability (in, out, undecided) of statements in an argument graph. The method is responsible for balancing conflicting pro and con arguments to select among competing options of each issue, handling any cycles in the argument graph and resolving attack relations among arguments.

**Argument graph**. A set of arguments

**Argument map**. A visualization of an argument graph used to explain and understand the arguments and their relationships.  Also called argument diagrams.

**Carneades**. An open source software system being developed by Fraunhofer FOKUS for supporting a variety of argumentation tasks, including argument construction, evaluation and visualization.  Arguments are constructed by applying rule-based models of argmentation schemes, laws and regulations to assumptions and facts, using an inference engine based on Constraint Handling Rules (CHR).  See <https://carneades.github.io/> for further information.

**Compliance checking**. A process of using an inference engine to apply a rule base modeling the GDPR to facts from a data use document to construct, evaluate and visualize arguments about whether or not the data use statements in the document comply with the regulation.  The resulting argument map is intended to be sufficient for explaining why the data uses statements are compliant or not, and for enabling authors to revise the data use statements in the document to achieve compliance, if needed.

**Constraint Handling Rules (CHR)**. A declarative, rule-based language, introduced in 1991 by Thom Frühwirth. Originally intended for constraint programming, CHR has found applications in abductive reasoning, multi-agent systems, natural language processing, compilation, scheduling, spatial-temporal reasoning, testing and verification, and type systems.  At Fraunhofer FOKUS, we have used CHR to represent defeasible inference rules (argumentation schemes) and model laws and regulations in the Carneades argumenation system.

**Data use document**. A sequence of data use statements, declared in a document.

**Data use statement**. A statement of the form "<data idenfication qualifier> <data category> from <scope source> is used by <use scope> to <action> the <result scope>", as defined in ISO/IEC ISO/IEC CD 19944.

**General Data Protection Regulation (GDPR)**. The new European data protection regulation, which is expected to be adopted in 2016 and will apply uniformly to all member states of the European Union.

**Inference Engine**. An inference engine applies rules in a rule base to a set of statements assumed to be true (called "assumptions" or "facts") to generate arguments.  Called "reasoning engine" in the contract.

**ISO/IEC CD 19944**. The ISO/IEC standard for "Information Technology - Cloud Computing - Data and their Flow across Devices and Cloud Services"

**Rule**. A representation, in a special purpose programming language, of a defeasible inference rule. Also called an "argumentation scheme".  There are many kind of rule-based programming languages, such as production rules, Constraint Handling Rules, and logic programming rules (as in Prolog), where the semantics of the rules differ considerably. Rules of the kind needed here, modeling defeasible inference rules, can be implemented via compilation to Constraint Handling Rules.

**Rule base**. A rule-based model of the some domain, consisting of a set of rules and metadata. In the context of DUCK, a rule-base will model the regulations of the GDPR and be used to check whether a data use document is compliant with the GDPR.

**Rule base editor**. A mode for customizing some programmer's text editor, such as Emacs, Vim, Atom or Visual Studio Code, to help modeller's to efficiently write and edit rule bases and detect and correct syntactical and semantic errors in rule bases.

**Statement**. A declarative sentence (proposition) which may be true or false. Data use statements are a special kind of statement.

# User Stories

## Compliance Engine User Stories

- As an author, I want to be able to check the compliance of data use documents with the European General Data Protection Regulation and receive a clear explanation of the ways in which the document is and is not compliant, sufficient to be able to locate the non-compliant statements and understand what changes are needed in order to achieve compliance.

- As an author, I want to check data use statements for compliance with data protection regulations, to comply with the law and minimize legal risks.

- As an author, I want to view diagrams of the arguments pro and con compliance, so as to have a better understanding of the compliance issues.

- As an author, I want recommendations for modifying my data use statements to resolve conflicts between my data use statements and the model of the data protection regulations.

## Rule Editor User Stories

- As a modeler, I want to check rule bases for syntactical and semantical errors and receive understandable error messages referencing particular points in the rule base, so as to be able to correct errors as soon as possible, before using the rule base with an inference engine to construct arguments.

- As a modeler, I want to model regulations using a programming language with plain text files, so as to be able to edit and exchange models with any text editor.

- As a modeler, I want to translate and view the rules of a model in natural language, to make it easier to validate the model in collaboration with lawyers.

# Tasks

The basic system architecture for meeting these requirements has already been defined.  See the architecture documents on GitHub.

Next steps for developing the compliance checker include (not necessarily in this order):

- Listing the syntactic and semantic errors to be identified by the rule editor.
- Defining the Go API (Interface) of the compliance checker
- Defining the RESTful API of the compliance checker, to be implemented as in integral part of the DUCK server backend.  (The compliance checker is a Go library which can be linked into the DUCK server, since the DUCK server is also being implemented in Go.
- Choosing the editor platform for the rule editor (e.g. Atom, Visual Studio Code or something else)
- Identify in which ways a data use document can be non-compliant with the GDPR.
- Modeling the parts of the GDPR relevant for checking the compliance of data use documents.
- Implementing the rule editor
- Implementing the compliance checker GO API
- Implementing the RESTful API of the compliance checker

To be able to demonstrate an initial version of the compliance checker as soon as possible, we propose the following procedure.

Torben, in collaboration with Tom, can

1. Define and implement the Go and RESTful APIs of the compliance checker
2. Integrate the compliance checker library into the DUCK server.
3. Implementing the Web user interfaces required to invoke the compliance checker.

In parallel, Tom can

1. Identify some initial relevant GDPR rules and model these in a rule base, without yet the benefit of the rule editor.  After this step, the basic idea of the compliance checker can be demonstrated and validated.
2. Design and implement the rule editor, with possibly Torben and/or Jim helping with the implementation.
3. Use the rule editor to refine and complete the GDPR rule base.
