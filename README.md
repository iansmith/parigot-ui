# parigot-ui

This is the place I'm doing a side project to make a UI layer for parigot.

There are three principal parts:

1. Using [fyne[(https://fyne.io) to build some desktop tools for parigot.  There is a an attempt a log viewer here already,
but it needs a lot of work.
2. A new toolkit called the wcl or Web Coordination Language.  The idea is to make many tasks of building websites 
language-neutral, with "holes" to drop into your primary programming language.  I'm not very happy with how complex it has gotten
and we probably need to remove some things.
3. The most important bit is the idea of `DOM Server` that is a parigot service that runs in a browser.  I built a proof of concept
to show the idea, but making it a full fledged tool is going to require a lot of work on the jascript side, plus the plubmbing to the 
golang or other pragramming language that you can use with the DOM service.  

Done correctly, number three would be awesomely cool as it would easily allow authors of back-end services to cerate front-end services 
using the same tooling and programming model they use for the back-end.








