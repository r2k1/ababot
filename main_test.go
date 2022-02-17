package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

var TestData = `[{"id":189786,"occurrenceId":195514,"resourceId":1,"start":"2022-02-13T06:00:00+13:00","end":"2022-02-13T07:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":189787,"occurrenceId":195515,"resourceId":3,"start":"2022-02-13T06:00:00+13:00","end":"2022-02-13T07:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":189791,"occurrenceId":195519,"resourceId":11,"start":"2022-02-13T06:00:00+13:00","end":"2022-02-13T07:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":189793,"occurrenceId":195521,"resourceId":2,"start":"2022-02-13T06:00:00+13:00","end":"2022-02-13T07:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":189794,"occurrenceId":195522,"resourceId":4,"start":"2022-02-13T06:00:00+13:00","end":"2022-02-13T07:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":189796,"occurrenceId":195524,"resourceId":8,"start":"2022-02-13T06:00:00+13:00","end":"2022-02-13T07:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":189797,"occurrenceId":195525,"resourceId":10,"start":"2022-02-13T06:00:00+13:00","end":"2022-02-13T07:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":194167,"occurrenceId":200732,"resourceId":5,"start":"2022-02-13T06:00:00+13:00","end":"2022-02-13T07:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":194168,"occurrenceId":200733,"resourceId":7,"start":"2022-02-13T06:00:00+13:00","end":"2022-02-13T07:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":197182,"occurrenceId":203867,"resourceId":6,"start":"2022-02-13T06:00:00+13:00","end":"2022-02-13T07:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":197184,"occurrenceId":203869,"resourceId":9,"start":"2022-02-13T06:00:00+13:00","end":"2022-02-13T09:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":197336,"occurrenceId":204039,"resourceId":12,"start":"2022-02-13T06:00:00+13:00","end":"2022-02-13T08:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":196333,"occurrenceId":202990,"resourceId":5,"start":"2022-02-13T07:00:00+13:00","end":"2022-02-13T09:00:00+13:00","title":"Phani Kowloori","rate":"34","status":"Confirmed","isCasual":false},{"id":196334,"occurrenceId":202991,"resourceId":10,"start":"2022-02-13T07:00:00+13:00","end":"2022-02-13T09:00:00+13:00","title":"Pankaj Chapke","rate":"30","status":"Confirmed","isCasual":false},{"id":196340,"occurrenceId":202997,"resourceId":11,"start":"2022-02-13T07:00:00+13:00","end":"2022-02-13T09:00:00+13:00","title":"Srivatsan Ganesh","rate":"36","status":"Confirmed","isCasual":false},{"id":197183,"occurrenceId":203868,"resourceId":6,"start":"2022-02-13T07:00:00+13:00","end":"2022-02-13T09:00:00+13:00","title":"Phani Kowloori","rate":"34","status":"Confirmed","isCasual":false},{"id":197752,"occurrenceId":204457,"resourceId":7,"start":"2022-02-13T07:00:00+13:00","end":"2022-02-13T08:00:00+13:00","title":"Dhanny Oud","rate":"15","status":"Confirmed","isCasual":false},{"id":197796,"occurrenceId":204501,"resourceId":3,"start":"2022-02-13T07:00:00+13:00","end":"2022-02-13T08:00:00+13:00","title":"Shumayala Syeda","rate":"44","status":"Confirmed","isCasual":false},{"id":179980,"occurrenceId":184715,"resourceId":4,"start":"2022-02-13T08:00:00+13:00","end":"2022-02-13T10:00:00+13:00","title":"Southern Shots (Sri Lanka Badminton Club)","rate":"34","status":"Confirmed","isCasual":false},{"id":179981,"occurrenceId":184716,"resourceId":3,"start":"2022-02-13T08:00:00+13:00","end":"2022-02-13T10:00:00+13:00","title":"Southern Shots (Sri Lanka Badminton Club)","rate":"34","status":"Confirmed","isCasual":false},{"id":179982,"occurrenceId":184717,"resourceId":2,"start":"2022-02-13T08:00:00+13:00","end":"2022-02-13T10:00:00+13:00","title":"Southern Shots (Sri Lanka Badminton Club)","rate":"34","status":"Confirmed","isCasual":false},{"id":179983,"occurrenceId":184718,"resourceId":1,"start":"2022-02-13T08:00:00+13:00","end":"2022-02-13T10:00:00+13:00","title":"Southern Shots (Sri Lanka Badminton Club)","rate":"34","status":"Confirmed","isCasual":false},{"id":196344,"occurrenceId":203001,"resourceId":8,"start":"2022-02-13T08:00:00+13:00","end":"2022-02-13T09:00:00+13:00","title":"Julie Ann Krishnakumar","rate":"15","status":"Confirmed","isCasual":false},{"id":196356,"occurrenceId":203013,"resourceId":7,"start":"2022-02-13T08:00:00+13:00","end":"2022-02-13T09:00:00+13:00","title":"Sharath Polapragada","rate":"15","status":"Confirmed","isCasual":false},{"id":197335,"occurrenceId":204038,"resourceId":12,"start":"2022-02-13T08:00:00+13:00","end":"2022-02-13T10:00:00+13:00","title":"Jon  Cook ","rate":"30","status":"Confirmed","isCasual":false},{"id":196282,"occurrenceId":202939,"resourceId":5,"start":"2022-02-13T09:00:00+13:00","end":"2022-02-13T10:00:00+13:00","title":"Jonathan Curtin","rate":"17","status":"Confirmed","isCasual":false},{"id":196345,"occurrenceId":203002,"resourceId":8,"start":"2022-02-13T09:00:00+13:00","end":"2022-02-13T10:00:00+13:00","title":"Julie Ann Krishnakumar","rate":"15","status":"Confirmed","isCasual":false},{"id":196355,"occurrenceId":203012,"resourceId":7,"start":"2022-02-13T09:00:00+13:00","end":"2022-02-13T10:00:00+13:00","title":"Sharath Polapragada","rate":"15","status":"Confirmed","isCasual":false},{"id":196359,"occurrenceId":203016,"resourceId":9,"start":"2022-02-13T09:00:00+13:00","end":"2022-02-13T10:00:00+13:00","title":"Arun Paluru","rate":"15","status":"Confirmed","isCasual":false},{"id":196889,"occurrenceId":203570,"resourceId":10,"start":"2022-02-13T09:00:00+13:00","end":"2022-02-13T10:00:00+13:00","title":"Nagabhushanam Gorantla","rate":"15","status":"Confirmed","isCasual":false},{"id":196925,"occurrenceId":203606,"resourceId":11,"start":"2022-02-13T09:00:00+13:00","end":"2022-02-13T10:00:00+13:00","title":"Anshul Somani","rate":"40","status":"Confirmed","isCasual":false},{"id":197354,"occurrenceId":204057,"resourceId":6,"start":"2022-02-13T09:00:00+13:00","end":"2022-02-13T10:00:00+13:00","title":"Jude Fernandes","rate":"17","status":"Confirmed","isCasual":false},{"id":181178,"occurrenceId":185932,"resourceId":1,"start":"2022-02-13T10:00:00+13:00","end":"2022-02-13T13:00:00+13:00","title":"Balmoral Badminton Club","rate":"51","status":"Confirmed","isCasual":false},{"id":181179,"occurrenceId":185933,"resourceId":2,"start":"2022-02-13T10:00:00+13:00","end":"2022-02-13T13:00:00+13:00","title":"Balmoral Badminton Club","rate":"51","status":"Confirmed","isCasual":false},{"id":181180,"occurrenceId":185934,"resourceId":3,"start":"2022-02-13T10:00:00+13:00","end":"2022-02-13T13:00:00+13:00","title":"Balmoral Badminton Club","rate":"51","status":"Confirmed","isCasual":false},{"id":181181,"occurrenceId":185935,"resourceId":5,"start":"2022-02-13T10:00:00+13:00","end":"2022-02-13T13:00:00+13:00","title":"Balmoral Badminton Club","rate":"51","status":"Confirmed","isCasual":false},{"id":181182,"occurrenceId":185936,"resourceId":4,"start":"2022-02-13T10:00:00+13:00","end":"2022-02-13T13:00:00+13:00","title":"Balmoral Badminton Club","rate":"51","status":"Confirmed","isCasual":false},{"id":181183,"occurrenceId":185937,"resourceId":6,"start":"2022-02-13T10:00:00+13:00","end":"2022-02-13T13:00:00+13:00","title":"Balmoral Badminton Club","rate":"51","status":"Confirmed","isCasual":false},{"id":181201,"occurrenceId":185955,"resourceId":7,"start":"2022-02-13T10:00:00+13:00","end":"2022-02-13T13:00:00+13:00","title":"Balmoral Badminton Club","rate":"45","status":"Confirmed","isCasual":false},{"id":181202,"occurrenceId":185956,"resourceId":9,"start":"2022-02-13T10:00:00+13:00","end":"2022-02-13T13:00:00+13:00","title":"Balmoral Badminton Club","rate":"45","status":"Confirmed","isCasual":false},{"id":181203,"occurrenceId":185957,"resourceId":8,"start":"2022-02-13T10:00:00+13:00","end":"2022-02-13T13:00:00+13:00","title":"Balmoral Badminton Club","rate":"45","status":"Confirmed","isCasual":false},{"id":181204,"occurrenceId":185958,"resourceId":10,"start":"2022-02-13T10:00:00+13:00","end":"2022-02-13T13:00:00+13:00","title":"Balmoral Badminton Club","rate":"45","status":"Confirmed","isCasual":false},{"id":181205,"occurrenceId":185959,"resourceId":11,"start":"2022-02-13T10:00:00+13:00","end":"2022-02-13T12:00:00+13:00","title":"Balmoral Badminton Club","rate":"30","status":"Confirmed","isCasual":false},{"id":197548,"occurrenceId":204253,"resourceId":12,"start":"2022-02-13T10:00:00+13:00","end":"2022-02-13T11:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":197550,"occurrenceId":204255,"resourceId":12,"start":"2022-02-13T11:00:00+13:00","end":"2022-02-13T12:00:00+13:00","title":"Warren Ji","rate":"15","status":"Confirmed","isCasual":false},{"id":197549,"occurrenceId":204254,"resourceId":12,"start":"2022-02-13T12:00:00+13:00","end":"2022-02-13T13:00:00+13:00","title":"Warren Ji","rate":"15","status":"Confirmed","isCasual":false},{"id":197792,"occurrenceId":204497,"resourceId":11,"start":"2022-02-13T12:00:00+13:00","end":"2022-02-13T13:00:00+13:00","title":"Hai Lan","rate":"18","status":"Confirmed","isCasual":false},{"id":193355,"occurrenceId":199724,"resourceId":7,"start":"2022-02-13T13:00:00+13:00","end":"2022-02-13T14:00:00+13:00","title":"Xu Li","rate":"14.5","status":"Confirmed","isCasual":false},{"id":197296,"occurrenceId":203999,"resourceId":11,"start":"2022-02-13T13:00:00+13:00","end":"2022-02-13T15:00:00+13:00","title":"Puntarika Meecharoen","rate":"80","status":"Confirmed","isCasual":true},{"id":197528,"occurrenceId":204233,"resourceId":4,"start":"2022-02-13T13:00:00+13:00","end":"2022-02-13T14:00:00+13:00","title":"Jerry Hu","rate":"17","status":"Confirmed","isCasual":false},{"id":197529,"occurrenceId":204234,"resourceId":1,"start":"2022-02-13T13:00:00+13:00","end":"2022-02-13T14:00:00+13:00","title":"Sharon Zhang","rate":"44","status":"Confirmed","isCasual":true},{"id":197530,"occurrenceId":204235,"resourceId":3,"start":"2022-02-13T13:00:00+13:00","end":"2022-02-13T14:00:00+13:00","title":"Tsuang  Hu ","rate":"17","status":"Confirmed","isCasual":false},{"id":197551,"occurrenceId":204256,"resourceId":12,"start":"2022-02-13T13:00:00+13:00","end":"2022-02-13T19:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":197776,"occurrenceId":204481,"resourceId":9,"start":"2022-02-13T13:00:00+13:00","end":"2022-02-13T14:00:00+13:00","title":"Tharindu Kaluarachchi","rate":"40","status":"Confirmed","isCasual":true},{"id":197779,"occurrenceId":204484,"resourceId":8,"start":"2022-02-13T13:00:00+13:00","end":"2022-02-13T14:00:00+13:00","title":"Zezheng DONG","rate":"40","status":"Confirmed","isCasual":true},{"id":197793,"occurrenceId":204498,"resourceId":10,"start":"2022-02-13T13:00:00+13:00","end":"2022-02-13T14:00:00+13:00","title":"Hai Lan","rate":"18","status":"Confirmed","isCasual":false},{"id":197813,"occurrenceId":204518,"resourceId":5,"start":"2022-02-13T13:00:00+13:00","end":"2022-02-13T14:00:00+13:00","title":"NATARAJ DEIVAMANI","rate":"22","status":"Confirmed","isCasual":false},{"id":197831,"occurrenceId":204536,"resourceId":6,"start":"2022-02-13T13:00:00+13:00","end":"2022-02-13T15:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":197848,"occurrenceId":204553,"resourceId":2,"start":"2022-02-13T13:00:00+13:00","end":"2022-02-13T14:00:00+13:00","title":"Eva Yin","rate":"44","status":"Confirmed","isCasual":true},{"id":193356,"occurrenceId":199725,"resourceId":7,"start":"2022-02-13T14:00:00+13:00","end":"2022-02-13T16:00:00+13:00","title":"Xu Li","rate":"29","status":"Confirmed","isCasual":false},{"id":196527,"occurrenceId":203184,"resourceId":8,"start":"2022-02-13T14:00:00+13:00","end":"2022-02-13T16:00:00+13:00","title":"David Xu","rate":"30","status":"Confirmed","isCasual":false},{"id":197042,"occurrenceId":203723,"resourceId":5,"start":"2022-02-13T14:00:00+13:00","end":"2022-02-13T16:00:00+13:00","title":"Raymond Biscocho","rate":"44","status":"Confirmed","isCasual":false},{"id":197043,"occurrenceId":203724,"resourceId":4,"start":"2022-02-13T14:00:00+13:00","end":"2022-02-13T16:00:00+13:00","title":"Tomas Morato","rate":"44","status":"Confirmed","isCasual":false},{"id":197113,"occurrenceId":203794,"resourceId":9,"start":"2022-02-13T14:00:00+13:00","end":"2022-02-13T15:00:00+13:00","title":"Cathy Yin","rate":"40","status":"Confirmed","isCasual":true},{"id":197213,"occurrenceId":203898,"resourceId":1,"start":"2022-02-13T14:00:00+13:00","end":"2022-02-13T16:00:00+13:00","title":"JIABAO WANG","rate":"34","status":"Confirmed","isCasual":false},{"id":197302,"occurrenceId":204005,"resourceId":2,"start":"2022-02-13T14:00:00+13:00","end":"2022-02-13T15:00:00+13:00","title":"Susan Vaz","rate":"44","status":"Confirmed","isCasual":true},{"id":197303,"occurrenceId":204006,"resourceId":10,"start":"2022-02-13T14:00:00+13:00","end":"2022-02-13T15:00:00+13:00","title":"Xuzhao Yan","rate":"18","status":"Confirmed","isCasual":false},{"id":197313,"occurrenceId":204016,"resourceId":3,"start":"2022-02-13T14:00:00+13:00","end":"2022-02-13T16:00:00+13:00","title":"Logan Burgess","rate":"34","status":"Confirmed","isCasual":false},{"id":195847,"occurrenceId":202504,"resourceId":11,"start":"2022-02-13T15:00:00+13:00","end":"2022-02-13T16:00:00+13:00","title":"Warren Ji","rate":"15","status":"Confirmed","isCasual":false},{"id":196336,"occurrenceId":202993,"resourceId":2,"start":"2022-02-13T15:00:00+13:00","end":"2022-02-13T17:00:00+13:00","title":"Enoch Wu","rate":"34","status":"Confirmed","isCasual":false},{"id":196346,"occurrenceId":203003,"resourceId":10,"start":"2022-02-13T15:00:00+13:00","end":"2022-02-13T17:00:00+13:00","title":"Xuzhao Yan","rate":"36","status":"Confirmed","isCasual":false},{"id":197094,"occurrenceId":203775,"resourceId":9,"start":"2022-02-13T15:00:00+13:00","end":"2022-02-13T17:00:00+13:00","title":"Jason Wong","rate":"30","status":"Confirmed","isCasual":false},{"id":197830,"occurrenceId":204535,"resourceId":6,"start":"2022-02-13T15:00:00+13:00","end":"2022-02-13T16:00:00+13:00","title":"Ben Yu","rate":"16.5","status":"Confirmed","isCasual":false},{"id":195848,"occurrenceId":202505,"resourceId":11,"start":"2022-02-13T16:00:00+13:00","end":"2022-02-13T17:00:00+13:00","title":"Warren Ji","rate":"15","status":"Confirmed","isCasual":false},{"id":196428,"occurrenceId":203085,"resourceId":3,"start":"2022-02-13T16:00:00+13:00","end":"2022-02-13T18:00:00+13:00","title":"Tony  Liu","rate":"34","status":"Confirmed","isCasual":false},{"id":196707,"occurrenceId":203376,"resourceId":4,"start":"2022-02-13T16:00:00+13:00","end":"2022-02-13T18:00:00+13:00","title":"Ben Yu","rate":"33","status":"Confirmed","isCasual":false},{"id":196987,"occurrenceId":203668,"resourceId":1,"start":"2022-02-13T16:00:00+13:00","end":"2022-02-13T18:00:00+13:00","title":"Sarah Park","rate":"44","status":"Confirmed","isCasual":false},{"id":197070,"occurrenceId":203751,"resourceId":5,"start":"2022-02-13T16:00:00+13:00","end":"2022-02-13T17:00:00+13:00","title":"mac ye","rate":"17","status":"Confirmed","isCasual":false},{"id":197379,"occurrenceId":204084,"resourceId":8,"start":"2022-02-13T16:00:00+13:00","end":"2022-02-13T17:00:00+13:00","title":"Lin Yang","rate":"40","status":"Confirmed","isCasual":true},{"id":197387,"occurrenceId":204092,"resourceId":7,"start":"2022-02-13T16:00:00+13:00","end":"2022-02-13T17:00:00+13:00","title":"Simon Li","rate":"15","status":"Confirmed","isCasual":false},{"id":197832,"occurrenceId":204537,"resourceId":6,"start":"2022-02-13T16:00:00+13:00","end":"2022-02-13T19:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":196387,"occurrenceId":203044,"resourceId":11,"start":"2022-02-13T17:00:00+13:00","end":"2022-02-13T19:00:00+13:00","title":"Yang Cai","rate":"36","status":"Confirmed","isCasual":false},{"id":196389,"occurrenceId":203046,"resourceId":10,"start":"2022-02-13T17:00:00+13:00","end":"2022-02-13T19:00:00+13:00","title":"Yuan Tsai","rate":"36","status":"Confirmed","isCasual":false},{"id":196418,"occurrenceId":203075,"resourceId":9,"start":"2022-02-13T17:00:00+13:00","end":"2022-02-13T19:00:00+13:00","title":"Hui Ying Khor","rate":"30","status":"Confirmed","isCasual":false},{"id":196419,"occurrenceId":203076,"resourceId":8,"start":"2022-02-13T17:00:00+13:00","end":"2022-02-13T19:00:00+13:00","title":"Ron Chan","rate":"30","status":"Confirmed","isCasual":false},{"id":196440,"occurrenceId":203097,"resourceId":5,"start":"2022-02-13T17:00:00+13:00","end":"2022-02-13T19:00:00+13:00","title":"Jimmy Lin","rate":"34","status":"Confirmed","isCasual":false},{"id":196984,"occurrenceId":203665,"resourceId":7,"start":"2022-02-13T17:00:00+13:00","end":"2022-02-13T18:00:00+13:00","title":"Warren Ji","rate":"15","status":"Confirmed","isCasual":false},{"id":197069,"occurrenceId":203750,"resourceId":2,"start":"2022-02-13T17:00:00+13:00","end":"2022-02-13T18:00:00+13:00","title":"mac ye","rate":"17","status":"Confirmed","isCasual":false},{"id":197035,"occurrenceId":203716,"resourceId":7,"start":"2022-02-13T18:00:00+13:00","end":"2022-02-13T19:00:00+13:00","title":"Josephine Lau","rate":"18","status":"Confirmed","isCasual":false},{"id":197811,"occurrenceId":204516,"resourceId":1,"start":"2022-02-13T18:00:00+13:00","end":"2022-02-13T19:00:00+13:00","title":"Tony Fang","rate":"17","status":"Confirmed","isCasual":false},{"id":197852,"occurrenceId":204557,"resourceId":4,"start":"2022-02-13T18:00:00+13:00","end":"2022-02-13T20:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":179079,"occurrenceId":183806,"resourceId":5,"start":"2022-02-13T19:00:00+13:00","end":"2022-02-13T21:00:00+13:00","title":"Friends United Badminton Club","rate":"34","status":"Confirmed","isCasual":false},{"id":179080,"occurrenceId":183807,"resourceId":6,"start":"2022-02-13T19:00:00+13:00","end":"2022-02-13T21:00:00+13:00","title":"Friends United Badminton Club","rate":"34","status":"Confirmed","isCasual":false},{"id":179081,"occurrenceId":192532,"resourceId":7,"start":"2022-02-13T19:00:00+13:00","end":"2022-02-13T21:00:00+13:00","title":"Friends United Badminton Club","rate":"30","status":"Confirmed","isCasual":false},{"id":179082,"occurrenceId":192533,"resourceId":8,"start":"2022-02-13T19:00:00+13:00","end":"2022-02-13T21:00:00+13:00","title":"Friends United Badminton Club","rate":"30","status":"Confirmed","isCasual":false},{"id":179083,"occurrenceId":192534,"resourceId":9,"start":"2022-02-13T19:00:00+13:00","end":"2022-02-13T21:00:00+13:00","title":"Friends United Badminton Club","rate":"30","status":"Confirmed","isCasual":false},{"id":179084,"occurrenceId":192535,"resourceId":10,"start":"2022-02-13T19:00:00+13:00","end":"2022-02-13T21:00:00+13:00","title":"Friends United Badminton Club","rate":"30","status":"Confirmed","isCasual":false},{"id":179086,"occurrenceId":192536,"resourceId":11,"start":"2022-02-13T19:00:00+13:00","end":"2022-02-13T21:00:00+13:00","title":"Friends United Badminton Club","rate":"30","status":"Confirmed","isCasual":false},{"id":179085,"occurrenceId":192537,"resourceId":12,"start":"2022-02-13T19:00:00+13:00","end":"2022-02-13T21:00:00+13:00","title":"Friends United Badminton Club","rate":"30","status":"Confirmed","isCasual":false},{"id":197380,"occurrenceId":204085,"resourceId":1,"start":"2022-02-13T19:00:00+13:00","end":"2022-02-13T20:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":194044,"occurrenceId":200609,"resourceId":4,"start":"2022-02-13T20:00:00+13:00","end":"2022-02-13T22:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":194045,"occurrenceId":200610,"resourceId":1,"start":"2022-02-13T20:00:00+13:00","end":"2022-02-13T22:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":197012,"occurrenceId":203693,"resourceId":3,"start":"2022-02-13T20:00:00+13:00","end":"2022-02-13T22:00:00+13:00","title":"Tim Yen","rate":"44","status":"Confirmed","isCasual":false},{"id":197285,"occurrenceId":203988,"resourceId":2,"start":"2022-02-13T20:00:00+13:00","end":"2022-02-13T21:00:00+13:00","title":"chao pang","rate":"22","status":"Confirmed","isCasual":false},{"id":197014,"occurrenceId":203695,"resourceId":2,"start":"2022-02-13T21:00:00+13:00","end":"2022-02-13T22:00:00+13:00","title":"Kelvin Choi","rate":"22","status":"Confirmed","isCasual":false},{"id":197023,"occurrenceId":203704,"resourceId":5,"start":"2022-02-13T21:00:00+13:00","end":"2022-02-13T22:00:00+13:00","title":"Bunnarath Chan","rate":"44","status":"Confirmed","isCasual":false},{"id":197025,"occurrenceId":203706,"resourceId":6,"start":"2022-02-13T21:00:00+13:00","end":"2022-02-13T22:00:00+13:00","title":"Leon Truong","rate":"44","status":"Confirmed","isCasual":true},{"id":197286,"occurrenceId":203989,"resourceId":7,"start":"2022-02-13T21:00:00+13:00","end":"2022-02-13T22:00:00+13:00","title":"chao pang","rate":"18","status":"Confirmed","isCasual":false},{"id":197290,"occurrenceId":203993,"resourceId":8,"start":"2022-02-13T21:00:00+13:00","end":"2022-02-13T22:00:00+13:00","title":"Everson  Zhong","rate":"18","status":"Confirmed","isCasual":false},{"id":197305,"occurrenceId":204008,"resourceId":12,"start":"2022-02-13T21:00:00+13:00","end":"2022-02-13T22:00:00+13:00","title":"Bonnie Lin","rate":"18","status":"Confirmed","isCasual":false},{"id":197507,"occurrenceId":204212,"resourceId":10,"start":"2022-02-13T21:00:00+13:00","end":"2022-02-13T22:00:00+13:00","title":"Zhaolun Miao","rate":"40","status":"Confirmed","isCasual":true},{"id":197804,"occurrenceId":204509,"resourceId":11,"start":"2022-02-13T21:00:00+13:00","end":"2022-02-13T22:00:00+13:00","title":"Aldric Khoo","rate":"15","status":"Confirmed","isCasual":false},{"id":189798,"occurrenceId":195526,"resourceId":1,"start":"2022-02-13T22:00:00+13:00","end":"2022-02-14T00:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":189799,"occurrenceId":195527,"resourceId":2,"start":"2022-02-13T22:00:00+13:00","end":"2022-02-14T00:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":189800,"occurrenceId":195528,"resourceId":4,"start":"2022-02-13T22:00:00+13:00","end":"2022-02-14T00:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":189801,"occurrenceId":195529,"resourceId":6,"start":"2022-02-13T22:00:00+13:00","end":"2022-02-14T00:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":189802,"occurrenceId":195530,"resourceId":8,"start":"2022-02-13T22:00:00+13:00","end":"2022-02-14T00:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":189803,"occurrenceId":195531,"resourceId":10,"start":"2022-02-13T22:00:00+13:00","end":"2022-02-14T00:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":189804,"occurrenceId":195532,"resourceId":12,"start":"2022-02-13T22:00:00+13:00","end":"2022-02-14T00:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":189805,"occurrenceId":195533,"resourceId":3,"start":"2022-02-13T22:00:00+13:00","end":"2022-02-14T00:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":189806,"occurrenceId":195534,"resourceId":5,"start":"2022-02-13T22:00:00+13:00","end":"2022-02-14T00:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":189807,"occurrenceId":195535,"resourceId":7,"start":"2022-02-13T22:00:00+13:00","end":"2022-02-14T00:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":189808,"occurrenceId":195536,"resourceId":9,"start":"2022-02-13T22:00:00+13:00","end":"2022-02-14T00:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false},{"id":189809,"occurrenceId":195537,"resourceId":11,"start":"2022-02-13T22:00:00+13:00","end":"2022-02-14T00:00:00+13:00","title":"Stadium and Operations Manager","rate":"0","status":"Confirmed","isCasual":false}]`

func testBookingData() []Booking {
	var result []Booking
	json.Unmarshal([]byte(TestData), &result)
	return result
}

func TestCalendar_Inverse(t *testing.T) {
	data, _ := fetchData()
	avaiable := availableSlots(data)
	fmt.Println(avaiable.toSlice())
}
