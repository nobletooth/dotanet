**DotaNet Ad Server**  
DotaNet Ad Server is a advertising platform that manages advertisers, publishers, and ad serving. It consists of several components including an ad server, event service, panel, and publisher website.

**Components**  
Ad Server
Event Service
Panel
Publisher Website

**Features**  
Ad retrieval and auction system
Click and impression tracking
Advertiser and publisher management
Real-time reporting
Support for multiple publisher websites

**Installation**  
Each component has its own build script. To build and deploy:

Navigate to the component directory
Run the build script:
```
./build.sh
```

**Structure**  
adserver/: Ad server component
eventservice/: Event tracking server
panel/: Management panel for advertisers and publishers
publisherwebsite/: Sample publisher websites
common/: Shared data structures
