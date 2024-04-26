# Workshop Title

Boss I crashed Production again! - How to prevent this with Chainsaw, Kyverno, and Keptn.

## Description

In this workshop, you will learn how to build and deploy applications with security, observeability, and reliability in mind. We will use Chainsaw to ensure your application behaves correctly, Kyverno to validate and enforce security policies in CI/CD Pipelines and Runtime. Futhermore you will learn how to make your deployment observeable with Keptn and how to automate the promotion of your application through different stages with a GitOps approach.

## How will your presentation benefit the ecosystem? 

This workshop will teach the attendees how to build a rock solid CI/CD Pipeline by combining multiple OpenSource Projects to ensure their applications are secure, observeable, and reliable.




## Agenda

1. Introduction and Workshop Goals (Charles & Christian) 10min
2. Setup Lab Environment (Christian) 15min
3. Make sth. Bad happened 45min
    TODO: What issues should we create?
     - Slow response time of the service
     - Distributed Tracing
     - Security Issue (PSP)
     - Helm Chart Issue (eg. secret name)
     - Root user in Container
     - Check for external dependencies payment service (external API available?)

4. Practice Session 20min?
    TODO: Think about possible Pitfalls.

TODO: Create sample application and chart


status:
  condition: bad


- deploy
    - post evaluation -> response time -> CR status succeeded -> Kyverno -> CR Analysis Resource -> Kyverno : YIPIII / NOOOO
    - post deployment tasl