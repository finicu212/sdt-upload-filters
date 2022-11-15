
Create a platform that allows to set up rules for ingesting and processing large amounts of data. The platform should allow selecting different data sources as streams (for example from network connections or files) and to configure a number of filters/handlers for each data stream. Result of each pipeline should also be collected in a configurable end destination. Filters can be duplicated, modified, composed. It must be possible to also attach loggers/observers to different parts of the pipeline for debugging.


Design patterns:

- Chain of Responsibility (data stream filters / handlers)
- Decorators (loggers / observers)
- Object Pool
- Builder / Factory for Connection OR Singleton Clipboard

![image](https://user-images.githubusercontent.com/44416281/201999065-1661b143-e832-4892-af23-0a3d6f560ba8.png)

Code https://pastebin.com/pE6UK4JB

![image](https://user-images.githubusercontent.com/44416281/201999113-634ced36-d705-4d67-adc4-5e46d9afda88.png)

Code https://pastebin.com/7az0ZqSr

![image](https://user-images.githubusercontent.com/44416281/201999139-5bb5d3c7-223c-4984-821d-fe88dd804168.png)

Code https://pastebin.com/NXWRBnBY


![image](https://user-images.githubusercontent.com/44416281/201999210-08f063e6-fe8e-486a-a71b-39b25a745c83.png)


Code https://pastebin.com/jfCfqir0

![image](https://user-images.githubusercontent.com/44416281/201999247-44882db7-c9a0-4069-92f7-0754aabf21ec.png)


https://pastebin.com/nXATujyU
