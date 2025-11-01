import http from "k6/http";
import { check } from "k6";

export const options = {
  // Define stages: 50 virtual users (VUs) ramp up over 10s,
  // stay at 50 VUs for 30s, then ramp down.
  stages: [
    { duration: "10s", target: 50 },
    { duration: "30s", target: 50 },
    { duration: "5s", target: 0 },
  ],
};

export default function () {
  // __VU is the virtual user ID
  // __ITER is the iteration number

  const payload = JSON.stringify({
    directory_path: "/var/data/loadtest",
    filename: `file-${__VU}-${__ITER}.txt`,
    file_type: "file",
    size: 12345,
    checksum: "checksum",
  });

  const params = {
    headers: {
      "Content-Type": "application/json",
    },
  };

  // Send the POST request
  const res = http.post("http://localhost:8000/files", payload, params);

  // Check if the request was successful (status code 201)
  check(res, {
    "created successfully (status 201)": (r) => r.status === 201,
  });
}
