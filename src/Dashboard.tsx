import React from "react";
import SqueueCard from "./board/SqueueCard";
import SinfoCard from "./board/SinfoCard";

// import cockpit from "cockpit";

const Dashboard = () => {
    return (
        <div id="dashboard">
            <SinfoCard />
            <SqueueCard />
        </div>
    );
};

export default Dashboard;
