export type PostConfig = {
  text: string;
  name: string;
  handle: string;
  verified: boolean;
  avatarUrl: string;
  date: string;
  likeCount: string;
  cta: string;
  width: "tight" | "wide";
  mode: "classic" | "simple";
};
