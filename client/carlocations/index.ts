import {CarLocationService} from "./services";
import {CarLocation} from "../../proto/car_location_pb";


 function listen(){
    const service = new CarLocationService("localhost:50051")
    const onReceiveCarLocation = (car: CarLocation) => {
        console.log("new car posted", car )
    }
    service.getAllCarLocationWithCallBack(onReceiveCarLocation)
        .then(closed =>  console.log("terminated", closed))
    .catch(err => console.log("error", err))
}

listen()