DUCK Architecture
===================

The DUCK architecture is separated into frontend and backend components. The frontend executes in an HTML5 browser (modern versions of Chrome, IE 10+, 
Firefox) while the backend is hosted in a server environment (Linux, Windows, OS X). The project will support backend deployment via a Docker container.

## Technology Summary

The backend is based on Golang. The frontend is based on Javascript, HTML 5, AngularJS and ZURB Foundation. In order to develop the backend, it is not 
necessary to be familiar with frontend technologies and vice versa. This allows DUCK to accommodate typical developer skillsets.    
   
## The Domain

DUCK is architected around a conceptual domain model based on the ISO 19944 Standard. This model is described in the _DUCK Domain_ slides. The principal 
entities in this model are:

- **User**: a person with an account in the system. A user may take on one or more roles as an _Author_ of a _Data Use Document_ or a modeller of a _Ruleset_.
- **Data Use Document**: a collection of ISO data use statements managed as a unit
- **Taxonomy**: an ISO-based vocabulary employed in a data use document to construct data use statements
- **Ruleset**: a collection of rules defining regulations enforced by a body such as a government. Data use documents are checked for compliance against a 
ruleset.
   
## Backend
   
The backend is written in Golang and is designed to be loosely coupled to clients. As such, it supports an HTTP RESTful API. This provides a documented and 
controlled method for accessing the backend by the default frontend UI as well by custom clients using it as a cloud service.  

The backend architecture is "process stateless" in that the only state maintained is stored directly in a database; the runtime process maintains no 
user-related state. In addition, the local file system is not used to store persistent information. This results in an architecture that is easy to scale 
horizontally, particularly in cloud and containerized environments where file systems and processes are ephemeral. In addition, this stateless architecture 
lends itself to more robust runtime operation: a backend process can be brought down and replaced without introducing a noticeable disruption to end-users.
 
## The REST resources
  
The REST API exposes the following resources. More detail can be found in the _DUCK Domain_ slides:
  
**User**
 
    /v1/users

Handles user management 

**Login**: 
    
    /v1/login
    
Creates a user login for authentication
     
**Document

    /v1/documents
    /v1/documents/{author id}/summary
    /v1/documents/{author id}/{doc id}

Retrieves, updates and deletes documents and document summaries (list view)

**Rule Set

    /v1/rulesets
    /v1/rulesets/{id}/
    /v1/rulesets/{id}/documents
    /v1/rulesets/{id}/documents/{id}
    
Retrieves, updates and deletes rulesets

The API is serviced by modular backend components managed by the Echo router. The router dispatches URI requests to an appropriate handler, which is 
responsible for executing application logic to process the request. This results in a modular design where backend code is organized by function. 

### Authentication and Identity Management   
   
Authentication and identity management is designed for high-availability, clustered environments and is based on [JSON Web Token](https://tools.ietf
.org/html/rfc7519). After a user logs in, an encrypted JWT is sent to the client which is set to expire after a configurable period of time. The token 
contains encrypted credentials. The client must present the token in an HTTP header for each request to the backend API. On receipt, the backend validates 
the token and determines whether to let the request process or reject it. The scheme allows the server to verify the validity of the token (relatively) 
cheaply without the need for persistent store lookups and to function easily in clustered and cloud environments where requests may be multiplexed to 
multiple backend runtimes. 

### Cluster Support   

Cluster support is at the forefront of the architecture design. Multiple backend runtimes can be clustered by interposing a standard HTTP load-balancer 
between clients and the runtimes. No other configuration is required.
 
### Compliance Checking

Compliance checking of data use documents against a ruleset is done using the Carnaedes argumentation engine. The engine is to be embedded (statically 
linked) as a Golang library. Embedding the engine (as opposed to accessing it over HTTP as a service) simplifies the end-user deployment experience as the 
backend will be a single executable that need only be configured to use a database.  

### Extensibility

The backend is designed to be extensible for key functionality. Specifically, it is possible for end-user developers to swap out the default database 
(CouchDB) and replace it with a custom implementation. Extensibility is provided by a plugin framework that enabled end-user developers to implement a 
defined interface (SPI) and link it to the backend during compilation. 
   
## Front End

The frontend is designed to be loosely coupled with the backend. It is based on two complimentary UI technologies:
    
* AngularJS
* ZURB Foundation
    
AngularJS provides UI-side service lifecycle management, databinding, templating, MVC, and REST communications with the backend. ZURB Foundation provides 
grid-based UI layout, cross-browser support, componentry, and typography. In addition, Sass (SCSS) is used for CSS extensibility.

The frontend UI is constrained to (mostly) coarse-grained requests to the backend as defined by the API. Consequently, nearly all rendering and validation 
will be performed on the client. This enhances fault-tolerance as well as performance since it avoids server invocation latency as well as introduces higher 
tolerance for network and server-side interruptions. From an end-user perspective, this will assist DUCK in providing a consistent and performant UI 
experience.  
 
    