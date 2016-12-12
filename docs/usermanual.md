## Account Creation

The first time after starting the DUCK application a new useraccount has to be created.


## Overwiev
After logging in the document overview  is shown. On this site existing documents can be accessed and new ones created.

Already populated testdocuments can be created by using the `Create Test Document` button on the upper right.
The  testdocuments are named after their validation result. As an example, the `No PII Example` testdocument will have no statement that relates to PII.

When hovering over a document, two actions become visible:
   - A `Edit` Button which opens the document view and
   - An icon of a bin which will delete that document

## Document view

In the document view two set of buttons are present. One on the left under the documents title, the other on the right under the validation result.
If a document has statements, these will be shown in the field under these buttons.

The left button set has, from left to right, a button to add a new statement, a button to save the document, one button to validate the document and one button to download the document.

The right set of buttons controls the documents language and the form of the statements, whether they are active or passive. The `Assumption Set` button has no functionality at this time.



### Statements

- each statement is shown in seperate box
- has ISO 19944 fields
- dropdown for each field
- Buttons:
    - Active/passive
    - cops statement
    - delete statement


### Validation
After validating the document via the `validate` button on the top left three new fields and a button are visible for each statement:

- `Relates To PII` or `Does not Relate to PII`: This field shows whether that statement relates to PII.
- `Legitimate Interest` or `No Legitimate Interest`: This field communicates whether that statement contains a legitimate interest of the CSP.
- `Consent Not Required` or `Consent Required`: This field shows whether this is a statement that requires consent.

The button `Show Compatible` reduces the statementsto only the ones that have a compatible purpose to that specific statement. To see all statements again the button `SHOW ALL` on the top of the page has to be used.

The documents *validation result* is displayed on the top right of the page.