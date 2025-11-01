import http from "k6/http";
import { check } from "k6";

export const options = {
  stages: [
    { duration: "10s", target: 50 }, // Ramp up to 50 users
    { duration: "30s", target: 50 }, // Stay at 50 users
    { duration: "5s", target: 0 }, // Ramp down
  ],
};

export default function () {
  const res = http.get("http://localhost:8000/files/2");

  check(res, {
    "retrieved successfully (status 200)": (r) => r.status === 200,
  });
}
