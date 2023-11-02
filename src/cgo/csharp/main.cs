using System.Runtime.InteropServices;

class Program
{
    [DllImport("ddns4cdn.dll", CharSet = CharSet.Unicode)]
    static extern void Ddns4cdnWorker([MarshalAs(UnmanagedType.LPStr)] string data);

    static void Main(string[] args)
    {
        string config = "config.toml";
        if (args.Length > 0)
        {
            config = args[0];
        }

        if (!File.Exists(config))
        {
            Console.Error.WriteLine($"Error opening file: {config}");
            return;
        }

        string data = File.ReadAllText(config);

        Ddns4cdnWorker(data);
    }
}

